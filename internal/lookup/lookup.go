package lookup

import (
	"net"

	"golang.org/x/xerrors"
)

type Func func(string) error

// HostnameByIP returns the hostname for the provided ip address.
func HostNameByIP(ip net.IP) (string, error) {
	hostnames, err := HostNamesByIP(ip)
	return hostnames[0], err
}

// HostnamesByIP returns all hostnames found for the provided ip address.
func HostNamesByIP(ip net.IP) ([]string, error) {
	hostnames, err := net.LookupAddr(ip.String())
	if err != nil {
		return nil, xerrors.Errorf("failed to lookup hostnames for ip address %q : %v", ip, err)
	}
	if len(hostnames) == 0 {
		return nil, xerrors.Errorf("no hostnames found for ip address: %q", ip)
	}
	return hostnames, nil
}

// AddrByHostName resolves the ip address of the provided hostname.
func AddrByHostName(hostname string) (*net.IP, error) {
	ipAddrs, err := AddrsByHostName(hostname)
	if err != nil {
		return nil, err
	}
	return ipAddrs[0], err
}

// AddrsByHostName returns all ip addresses found for the provided hostname.
func AddrsByHostName(hostname string) ([]*net.IP, error) {
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		return nil, xerrors.Errorf("failed to look up ip addresses for hostname %q: %v", hostname, err)
	}

	if len(addrs) == 0 {
		return nil, xerrors.Errorf("no ip addresses found for hostname %q", hostname)
	}

	var ipAddrs []*net.IP

	for _, a := range addrs {
		ipAddr := net.ParseIP(a)
		if ipAddr == nil {
			continue
		}
		ipAddrs = append(ipAddrs, &ipAddr)
	}
	return ipAddrs, nil
}

// NameServersByHostName looks up all nameservers for the provided hostname.
func NameServersByHostName(hostname string) ([]*net.NS, error) {
	nameServers, err := net.LookupNS(hostname)
	if err != nil {
		return nil, xerrors.Errorf("failed to look up name server for hostname %q : %v", hostname, err)
	}
	if len(nameServers) == 0 {
		return nil, xerrors.Errorf("no name servers found for hostname %q", hostname)
	}
	return nameServers, nil
}

// NetworkByHost returns the network address for the provided hostname.
func NetworkByHost(host string) (*net.IPMask, error) {
	ipAddr := net.ParseIP(host)
	if ipAddr == nil {
		addr, err := AddrByHostName(host)
		if err != nil {
			return nil, xerrors.Errorf("%q is an invalild host: %v", host, err)
		}
		ipAddr = *addr
	}

	network := ipAddr.DefaultMask()
	if network == nil {
		return nil, xerrors.Errorf("failed to get network address of host %q", ipAddr.String())
	}
	return &network, nil
}
