package scan

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	u "github.com/fuskovic/networker/internal/utils"
)

// Scanner contains fields that tell it what kind of scan needs to be performed.
type Scanner struct {
	Host    string
	Ports   []int
	Verbose bool
}

// Scan iterates through `portsForScanning` and scans
// each port in it's own goroutine. Scan does not exit
// until all goroutines are done.
func (s *Scanner) Scan(ctx context.Context) {
	var (
		wg   sync.WaitGroup
		open []string
		log  = sloghuman.Make(os.Stdout)
	)

	log.Info(ctx, "target", s.fields()...)
	ch := make(chan result)

	go func() {
		for r := range ch {
			if r.open {
				open = append(open, fmt.Sprintf("%d", r.port))
			}
			if s.Verbose {
				log.Info(ctx, u.TCP, r.fields()...)
			}
		}
	}()

	wg.Add(len(s.Ports))

	for _, port := range s.Ports {
		go func(p int) {
			s.scan(p, ch)
			wg.Done()
		}(port)
	}
	wg.Wait()
	close(ch)

	f := slog.F("open-ports", strings.Join(open, ","))
	log.Info(ctx, "scan complete", f)
}

func (s *Scanner) scan(port int, c chan<- result) {
	check := func(protocol string) {
		c <- result{
			addr:     s.Host,
			port:     port,
			protocol: protocol,
			open:     isOpen(protocol, s.Host, port),
		}
	}

	if s.ready(u.TCP, port) {
		check(u.TCP)
	}
}

func (s *Scanner) fields() []slog.Field {
	var fields []slog.Field

	ip := net.ParseIP(s.Host)
	if ip != nil {
		fields = append(fields,
			slog.F("hostname", u.HostNameByIP(s.Host)),
			slog.F("addr", s.Host),
		)
	} else {
		fields = append(fields,
			slog.F("hostname", s.Host),
			slog.F("addr", u.AddrByHostName(s.Host)))
	}
	return fields
}

func (s *Scanner) ready(protocol string, port int) bool {
	return s.Verbose && isOpen(protocol, s.Host, port) || !s.Verbose
}

func isOpen(protocol, host string, port int) bool {
	t := 3 * time.Second
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout(protocol, addr, t)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
