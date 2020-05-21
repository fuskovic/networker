package scan

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// TotalPorts is the total number of all tcp/udp ports
	TotalPorts = 65535
	tcp        = "tcp"
	udp        = "udp"
)

var (
	timeOut = 3 * time.Second
	stars   = strings.Repeat("*", 30)
)

type (
	// Config collects the command parameters for the scan sub-command.
	Config struct {
		IP                         string
		Ports                      []int
		UpTo                       int
		TCPonly, UDPonly, OpenOnly bool
	}
	scanner struct {
		host                       string
		tcpOnly, udpOnly, openOnly bool
	}
)

// Run executes the command logic for this package.
func Run(cfg *Config) error {
	if net.ParseIP(cfg.IP) == nil {
		return fmt.Errorf("%s is not a valid IP address", cfg.IP)
	}

	scanner := newScanner(cfg)

	switch {
	case cfg.UpTo > TotalPorts:
		return fmt.Errorf("can not scan more than %d ports", TotalPorts)
	case cfg.UpTo > TotalPorts:
		return fmt.Errorf("can not scan more than %d ports", TotalPorts)
	case len(cfg.Ports) > 0:
		scanner.scanPorts(cfg.Ports)
	case cfg.UpTo > 0:
		scanner.scanUpTo(cfg.UpTo)
	default:
		scanner.scanAllPorts()
	}
	log.Println("scan complete")
	return nil
}

func newScanner(cfg *Config) *scanner {
	if !protocolSpecified(cfg.TCPonly, cfg.UDPonly) {
		log.Println("protocol unspecified enabling scanner for both")
		cfg.TCPonly = true
		cfg.UDPonly = true
	}

	return &scanner{
		host:     cfg.IP,
		tcpOnly:  cfg.TCPonly,
		udpOnly:  cfg.UDPonly,
		openOnly: cfg.OpenOnly,
	}
}

func (s *scanner) scanPorts(specifiedPorts []int) {
	s.start(specifiedPorts)
}

func (s *scanner) scanUpTo(upTo int) {
	portsForScanning := portsToScan(upTo)
	s.start(portsForScanning)
}

func (s *scanner) scanAllPorts() {
	portsForScanning := portsToScan(TotalPorts)
	s.start(portsForScanning)
}

func (s *scanner) scan(port int, c chan<- string) {
	if s.tcpOnly {
		if s.shouldLog(tcp, port) {
			c <- fmt.Sprintf("%s\nport : %s\nOpen : %t",
				stars,
				fmt.Sprintf("%s/%d", tcp, port),
				isOpen(tcp, s.host, port),
			)
		}
	}

	if s.udpOnly {
		if s.shouldLog(udp, port) {
			c <- fmt.Sprintf("%s\nport : %s\nOpen : %t",
				stars,
				fmt.Sprintf("%s/%d", udp, port),
				isOpen(udp, s.host, port),
			)
		}
	}
}

func (s *scanner) shouldLog(protocol string, port int) bool {
	return s.openOnly && isOpen(protocol, s.host, port) || !s.openOnly
}

func isOpen(protocol, host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout(protocol, address, timeOut)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func portsToScan(max int) (ports []int) {
	for port := 0; port < max; port++ {
		ports = append(ports, port)
	}
	return
}

func protocolSpecified(tcp, udp bool) bool {
	return tcp == true || udp == true
}

func organize(results []string) []string {
	var organized []string
	for i := 0; i < len(results); i++ {
		for _, r := range results {
			if strings.Contains(r, strconv.Itoa(i)) {
				organized = append(organized, r)
			}
		}
	}
	return organized
}

func (s *scanner) start(portsForScanning []int) {
	var wg sync.WaitGroup
	var results []string
	wg.Add(len(portsForScanning))
	ch := make(chan string)

	go func() {
		for {
			select {
			case r, ok := <-ch:
				results = append(results, r)
				if !ok {
					return
				}
			}
		}
	}()

	log.Println("starting scan...")

	for _, port := range portsForScanning {
		go func(p int) {
			s.scan(p, ch)
			wg.Done()
		}(port)
	}
	wg.Wait()
	close(ch)

	for _, r := range organize(results) {
		fmt.Println(r)
	}
}
