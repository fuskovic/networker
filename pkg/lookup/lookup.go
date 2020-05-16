package lookup

import (
	"fmt"
	"net"
	"strings"
)

// HostNamesByIP looks up all hostnames for an IP address.
func HostNamesByIP(IP string) error {
	IP = trim(IP)

	if net.ParseIP(IP) == nil {
		return fmt.Errorf("%s is not a valid IP address", IP)
	}

	fmt.Printf("Looking up hostnames for IP address: %s\n", IP)
	hostnames, err := net.LookupAddr(IP)
	if err != nil {
		return fmt.Errorf("failed to lookup hostnames for %s\nerror : %v", IP, err)
	}

	if len(hostnames) == 0 {
		return fmt.Errorf("no hostnames found")
	}

	for _, hn := range hostnames {
		fmt.Println(hn)
	}
	return nil
}

// IPsByHostName looks up all IP addresses for a hostname.
func IPsByHostName(hostName string) error {
	hostName = trim(hostName)

	fmt.Printf("Looking up IP addresses for hostname: %s\n", hostName)
	IPs, err := net.LookupHost(hostName)
	if err != nil {
		return fmt.Errorf("failed to look up IP addresses for %s\nerror : %v", hostName, err)
	}

	if len(IPs) == 0 {
		return fmt.Errorf("no IP addresses found")
	}

	for _, IP := range IPs {
		fmt.Println(IP)
	}
	return nil
}

// NameServersByHostName looks up all name servers for a hostname.
func NameServersByHostName(hostName string) error {
	hostName = trim(hostName)

	fmt.Printf("Looking up nameservers for %s\n", hostName)
	nameservers, err := net.LookupNS(hostName)
	if err != nil {
		return fmt.Errorf("failed to look up name server for %s\nerror : %v", hostName, err)
	}

	if len(nameservers) == 0 {
		return fmt.Errorf("no name servers found")
	}

	for _, ns := range nameservers {
		fmt.Println(ns.Host)
	}
	return nil
}

// NetworkByHostName looks up the network for a hostname.
func NetworkByHostName(hostName string) error {
	ip, err := net.ResolveIPAddr("ip", hostName)
	if err != nil {
		return fmt.Errorf("failed to resolve IP address from hostname : %s\nerror : %v", hostName, err)
	}

	addr := net.ParseIP(ip.String())
	if addr == nil {
		return fmt.Errorf("failed to validate the resolved IP : %s for hostname : %s", addr, hostName)
	}

	mask := addr.DefaultMask()
	network := addr.Mask(mask)
	format := "Hostname : %s\nAddress : %s\nNetwork : %s\n"
	fmt.Printf(format, hostName, addr.String(), network)
	return nil
}

func trim(s string) string {
	return strings.TrimSpace(s)
}
