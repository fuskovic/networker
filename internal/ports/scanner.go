package ports

import (
	"context"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	wellKnownPorts = 1024
	allPorts       = 65535
)

type Scanner interface {
	Scan(context.Context) map[string][]int
}

type scanner struct {
	sync.Mutex
	scans         map[string][]int
	shouldScanAll bool
	hosts         []string
}

// NewScanner initializes a new port-scanner based on whether or not the user wants to scan all ports or just the well-known ports.
func NewScanner(hosts []string, shouldScanAll bool) Scanner {
	scans := make(map[string][]int)
	for _, host := range hosts {
		scans[host] = []int{}
	}

	return &scanner{
		Mutex:         sync.Mutex{},
		scans:         scans,
		shouldScanAll: shouldScanAll,
	}
}

func (s *scanner) Scan(ctx context.Context) map[string][]int {
	var wg sync.WaitGroup
	for host := range s.scans {
		wg.Add(1)
		go func(h string) {
			s.scanHost(h)
			wg.Done()
		}(host)
	}
	wg.Wait()
	return s.scans
}

func (s *scanner) scanHost(host string) {
	var wg sync.WaitGroup
	for _, port := range portsToScan(s.shouldScanAll) {
		wg.Add(1)
		go func(p int) {
			if isOpen(host, p) {
				s.add(host, p)
			}
			wg.Done()
		}(port)
	}
	wg.Wait()
}

func (s *scanner) add(host string, port int) {
	s.Lock()
	s.scans[host] = append(s.scans[host], port)
	s.Unlock()
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
