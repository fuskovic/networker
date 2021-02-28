package ports

import (
	"context"
	"net"
	"strconv"
	"sync"
	"time"
)

type Scanner interface {
	Scan(context.Context, []int) []int
}

type scanner struct {
	sync.Mutex
	openPorts []int
	host      string
}

func NewScanner(host string) Scanner {
	return &scanner{
		Mutex: sync.Mutex{},
		host:  host,
	}
}

func (s *scanner) Scan(ctx context.Context, portsToScan []int) []int {
	var wg sync.WaitGroup
	for _, port := range portsToScan {
		wg.Add(1)
		go func(p int) {
			if isOpen(s.host, p) {
				s.add(p)
			}
			wg.Done()
		}(port)
	}
	wg.Wait()
	return s.openPorts
}

func (s *scanner) add(port int) {
	s.Lock()
	defer s.Unlock()
	s.openPorts = append(s.openPorts, port)
}

func isOpen(host string, port int) bool {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
