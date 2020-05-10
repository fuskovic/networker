package scan

import (
	"fmt"
	"log"
	"net"
	"time"
)

var timeOut = 10 * time.Second

// Scanner describes the port scanner.
type Scanner struct {
	Host             string
	TCPOnly, UDPOnly bool
}

// NewScanner initializes a new scanner.
func NewScanner(host string, tcpOnly, udpOnly bool) *Scanner {
	if !protocolSpecified(tcpOnly, udpOnly) {
		log.Println("protocol unspecified enabling scanner for both")
		tcpOnly = true
		udpOnly = true
	}

	return &Scanner{
		Host:    host,
		TCPOnly: tcpOnly,
		UDPOnly: udpOnly,
	}
}

// ScanPorts scans an explicit set of ports specified.
func (s *Scanner) ScanPorts(ports []int) {
	// TODO : implement
}

// ScanUpTo scans all ports up to the value specified.
func (s *Scanner) ScanUpTo(upTo int) {
	// TODO : implement
}

// ScanAllPorts scans all ports.
func (s *Scanner) ScanAllPorts() {
	// TODO : implement
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

func protocolSpecified(tcp, udp bool) bool {
	return tcp == true || udp == true
}
