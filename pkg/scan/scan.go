package scan

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	// TotalPorts is the total number of all tcp/udp ports
	TotalPorts = 65535
	tcp        = "tcp"
	udp        = "udp"
)

var timeOut = 3 * time.Second

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
	var wg sync.WaitGroup
	wg.Add(len(specifiedPorts))

	for _, port := range specifiedPorts {
		go func(p int) {
			s.scan(p)
			wg.Done()
		}(port)
	}
	wg.Wait()
}

// ScanUpTo scans all ports up to the value specified.
func (s *Scanner) ScanUpTo(upTo int) {
	var wg sync.WaitGroup
	portsForScanning := portsToScan(upTo)
	wg.Add(len(portsForScanning))

	for _, port := range portsForScanning {
		go func(p int) {
			s.scan(p)
			wg.Done()
		}(port)
	}
	wg.Wait()
}

// ScanAllPorts scans all ports.
func (s *Scanner) ScanAllPorts() {
	var wg sync.WaitGroup
	portsForScanning := portsToScan(TotalPorts)
	wg.Add(len(portsForScanning))

	for _, port := range portsToScan(TotalPorts) {
		go func(p int) {
			s.scan(p)
			wg.Done()
		}(port)
	}
	wg.Wait()
}

func (s *Scanner) scan(port int) {
	if s.tcpOnly {
		if s.shouldLog(tcp, port) {
			log.Printf("port : %s Open : %t\n",
				fmt.Sprintf("%s/%d", tcp, port),
				isOpen(tcp, s.host, port),
			)
		}
	}

	if s.udpOnly {
		if s.shouldLog(udp, port) {
			log.Printf("port : %s Open : %t\n",
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
