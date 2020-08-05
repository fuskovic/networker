package utils

import (
	"net"
	"os"
	"strings"
	"syscall"

	"cdr.dev/slog"
)

const (
	TotalPorts = 65536
	TCP        = "tcp"
	UDP        = "udp"
	Unknown    = "unknown"
)

var (
	Methods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	Signals = []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}
)

// Row contains the fields of a structured log.
type Row []slog.Field

// Valid checks if the protocol field is unknown or not.
func (r *Row) Valid() bool {
	for _, f := range *r {
		if f.Name == "proto" {
			if f.Value == Unknown {
				return false
			}
			break
		}
	}
	return true
}

// Add adds a field to the current row
func (r *Row) Add(name string, val interface{}) {
	*r = append(*r, slog.F(name, val))
}

// HostnameByIP prints any hostnames found for a given IP.
func HostNameByIP(a string) string {
	if net.ParseIP(a) == nil {
		return Unknown
	}

	hostnames, err := net.LookupAddr(a)
	if err != nil {
		return Unknown
	}

	if len(hostnames) == 0 {
		return Unknown
	}
	return hostnames[0]
}

// AddrByHostName prints any addrs found for a given hostname.
func AddrByHostName(hn string) string {
	addrs, err := net.LookupHost(hn)
	if err != nil {
		return Unknown
	}

	if len(addrs) == 0 {
		return Unknown
	}
	return addrs[0]
}

// HasProtoScheme evaluates whether or not
// 'http://' or 'https://' is present.
func HasProtoScheme(a string) bool {
	has := func(s string) bool {
		return strings.Contains(a, s)
	}
	return has("http://") || has("https://")
}
