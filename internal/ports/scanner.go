package ports

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/fuskovic/networker/internal/resolve"
)

var (
	wellKnownPorts = 1024
	allPorts       = 65535
)

type Scan struct {
	IP    string `json:"ip" table:"IP"`
	Host  string `json:"hostname" table:"Hostname"`
	Ports []int  `json:"open_ports" table:"OpenPorts"`
}

type Scanner interface {
	Scan(context.Context) ([]Scan, error)
}

type scanner struct {
	sync.Mutex
	scans         map[string][]int
	shouldScanAll bool
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

func (s *scanner) Scan(ctx context.Context) ([]Scan, error) {
	var wg sync.WaitGroup
	for host := range s.scans {
		wg.Add(1)
		go func(ip string) {
			s.scanHost(ip)
			wg.Done()
		}(host)
	}
	wg.Wait()

	var scans []Scan
	for ip, ports := range s.scans {
		hostname, _, err := resolve.HostAndAddr(ip)
		if err != nil {
			return nil, fmt.Errorf("failed to lookup hostname by ip for %s: %w", ip, err)
		}
		scans = append(scans, Scan{
			IP:    ip,
			Host:  hostname,
			Ports: ports,
		})
	}
	return scans, nil
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

func (s *scanner) add(ip string, port int) {
	s.Lock()
	s.scans[ip] = append(s.scans[ip], port)
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
