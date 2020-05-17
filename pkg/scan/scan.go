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

// Scanner describes the port scanner.
type Scanner struct {
	host                       string
	tcpOnly, udpOnly, openOnly bool
}

// NewScanner initializes a new scanner.
func NewScanner(host string, tcpOnly, udpOnly, openOnly bool) *Scanner {
	if !protocolSpecified(tcpOnly, udpOnly) {
		log.Println("protocol unspecified enabling scanner for both")
		tcpOnly = true
		udpOnly = true
	}

	return &Scanner{
		host:     host,
		tcpOnly:  tcpOnly,
		udpOnly:  udpOnly,
		openOnly: openOnly,
	}
}

// ScanPorts scans an explicit set of ports specified.
func (s *Scanner) ScanPorts(specifiedPorts []int) {
	s.start(specifiedPorts)
}

// ScanUpTo scans all ports up to the value specified.
func (s *Scanner) ScanUpTo(upTo int) {
	portsForScanning := portsToScan(upTo)
	s.start(portsForScanning)
}

// ScanAllPorts scans all ports.
func (s *Scanner) ScanAllPorts() {
	portsForScanning := portsToScan(TotalPorts)
	s.start(portsForScanning)
}

func (s *Scanner) scan(port int, c chan<- string) {
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

func (s *Scanner) shouldLog(protocol string, port int) bool {
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

func (s *Scanner) start(portsForScanning []int) {
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
