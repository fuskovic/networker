package ports

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	goping "github.com/tatsushid/go-fastping"

	"github.com/fuskovic/networker/internal/resolve"
)

var (
	wellKnownPorts = 1024
	allPorts       = 65535
)

type Scan struct {
	IP    string `json:"ip" table:"IP"`
	Host  string `json:"hostname" table:"HOSTNAME"`
	Ports []int  `json:"open_ports" table:"OPEN_PORTS"`
	Up    bool   `json:"up" yaml:"up" table:"UP"`
}

type Scanner interface {
	Scan(context.Context) ([]Scan, error)
}

type scanner struct {
	sync.Mutex
	scans         []Scan
	shouldScanAll bool
}

// NewScanner initializes a new port-scanner based on whether or not the user wants to scan all ports or just the well-known ports.
func NewScanner(hosts []string, shouldScanAll bool) Scanner {
	var (
		scans []Scan
		wg    sync.WaitGroup
		mu    sync.Mutex
	)

	for _, host := range hosts {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			p := goping.NewPinger()
			_, _ = p.Network("udp")

			netProto := "ip4:icmp"
			if strings.Index(ip, ":") != -1 {
				netProto = "ip6:ipv6-icmp"
			}

			addr, err := net.ResolveIPAddr(netProto, ip)
			if err != nil {
				return
			}

			p.AddIPAddr(addr)
			p.MaxRTT = time.Second

			s := Scan{IP: ip}
			p.OnRecv = func(addr *net.IPAddr, t time.Duration) { s.Up = true }
			if err := p.Run(); err != nil {
				return
			}

			mu.Lock()
			scans = append(scans, s)
			mu.Unlock()
		}(host)
	}
	wg.Wait()

	return &scanner{
		Mutex:         sync.Mutex{},
		scans:         scans,
		shouldScanAll: shouldScanAll,
	}
}

func (s *scanner) Scan(ctx context.Context) ([]Scan, error) {
	var wg sync.WaitGroup
	for _, scan := range s.scans {
		if scan.Up {
			wg.Add(1)
			go func(ip string) {
				defer wg.Done()
				s.scanHost(ip)
			}(scan.IP)
		}
	}
	wg.Wait()

	for i := range s.scans {
		hostname, _, err := resolve.HostAndAddr(s.scans[i].IP)
		if err != nil {
			return nil, fmt.Errorf("failed to lookup hostname by ip for %s: %w", s.scans[i].IP, err)
		}
		s.scans[i].Host = hostname
	}
	return s.scans, nil
}

func (s *scanner) scanHost(host string) {
	var wg sync.WaitGroup
	for _, port := range portsToScan(s.shouldScanAll) {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			if isOpen(host, p) {
				s.add(host, p)
			}
		}(port)
	}
	wg.Wait()
}

func (s *scanner) add(ip string, port int) {
	s.Lock()
	for i := range s.scans {
		if s.scans[i].IP == ip {
			s.scans[i].Ports = append(s.scans[i].Ports, port)
		}
	}
	s.Unlock()
}

func isOpen(ip string, port int) bool {
	addr := net.JoinHostPort(ip, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func portsToScan(shouldScanAll bool) []int {
	var max int
	if shouldScanAll {
		max = allPorts
	} else {
		max = wellKnownPorts
	}

	var ports []int
	for p := 0; p < max; p++ {
		ports = append(ports, p)
	}
	return ports
}
