package internal

import (
	"fmt"
	"net"
	"strings"

	"go.coder.com/flog"
)

// LookUpFunc is a type that prints look-up results.
type LookUpFunc func(string) error

// HostnamesByIP prints any hostnames found for a given IP.
func HostNamesByIP(a string) error {
	a = trim(a)
	flog.Info("looking up hostnames for %s", a)

	if net.ParseIP(a) == nil {
		return fmt.Errorf("%s is not a valid IP address", a)
	}

	hostnames, err := net.LookupAddr(a)
	if err != nil {
		return fmt.Errorf("failed to lookup hostnames for %s\nerror : %v", a, err)
	}

	if len(hostnames) == 0 {
		return fmt.Errorf("no hostnames found")
	}

	for _, hn := range hostnames {
		flog.Info(hn)
	}
	return nil
}

// AddrsByHostName prints any addrs found for a given hostname.
func AddrsByHostName(hn string) error {
	hn = trim(hn)
	flog.Info("looking up addresses for %s", hn)

	addrs, err := net.LookupHost(hn)
	if err != nil {
		return fmt.Errorf("failed to look up IP addresses for %s\nerror : %v", hn, err)
	}

	if len(addrs) == 0 {
		return fmt.Errorf("no IP addresses found")
	}

	for _, a := range addrs {
		flog.Info(a)
	}
	return nil
}

// NameServersByHostName prints any name servers found for a given hostname.
func NameServersByHostName(hn string) error {
	hn = trim(hn)
	flog.Info("looking up name servers for %s", hn)

	nameServers, err := net.LookupNS(hn)
	if err != nil {
		return fmt.Errorf("failed to look up name server for %s\nerror : %v", hn, err)
	}

	if len(nameServers) == 0 {
		return fmt.Errorf("no name servers found")
	}

	for _, ns := range nameServers {
		flog.Info(ns.Host)
	}
	return nil
}

// NwByHostName prints a network found for a given hostname.
func NwByHostName(hn string) error {
	ip, err := net.ResolveIPAddr("ip", hn)
	if err != nil {
		return fmt.Errorf("failed to resolve IP address from hostname : %s\nerror : %v", hn, err)
	}

	a := net.ParseIP(ip.String())
	if a == nil {
		return fmt.Errorf("failed to validate the resolved IP : %s for hostname : %s", a, hn)
	}
	m := a.DefaultMask()
	nw := a.Mask(m)
	flog.Info("network : %s", nw)
	return nil
}

func trim(s string) string {
	return strings.TrimSpace(s)
}
