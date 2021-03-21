package resolve

import (
	"net"
	"time"

	"github.com/ammario/ipisp"
	"golang.org/x/xerrors"
)

type Func func(string) error

// InternetServiceProvider describes an internet service provider.
type InternetServiceProvider struct {
	Name                    string     `json:"name" table:"Name"`
	IP                      *net.IP    `json:"ip_address" table:"IP"`
	Country                 string     `json:"country" table:"Country"`
	Registry                string     `json:"registry" table:"Registry"`
	IpRange                 *net.IPNet `json:"ip_range" table:"IP-Range"`
	AutonomousServiceNumber string     `json:"autonomous_service_number" table:"ASN"`
	AllocatedAt             *time.Time `json:"allocated_at" table:"AllocatedAt"`
}

// NameServer is used in place of the standard library object to support table writes.
type NameServer struct {
	Host string `json:"host" table:"Host"`
}

// HostNameByIP returns the hostname for the provided ip address.
func HostNameByIP(ip net.IP) (string, error) {
	hostnames, err := HostNamesByIP(ip)
	if err != nil {
		return "", err
	}
	return hostnames[0], nil
}

// HostNamesByIP returns all hostnames found for the provided ip address.
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
	if len(ipAddrs) == 0 {
		return nil, xerrors.Errorf("no addresses found for hostname %q", hostname)
	}
	ipv4 := ipAddrs[0].To4()
	if ipv4 == nil {
		return nil, xerrors.Errorf("failed to cast %q to ipv4", ipAddrs[0])
	}
	return &ipv4, err
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
func NameServersByHostName(hostname string) ([]NameServer, error) {
	internalNameServers, err := net.LookupNS(hostname)
	if err != nil {
		return nil, xerrors.Errorf("failed to look up name server for hostname %q : %v", hostname, err)
	}
	if len(internalNameServers) == 0 {
		return nil, xerrors.Errorf("no name servers found for hostname %q", hostname)
	}
	var nameServers []NameServer
	for _, ns := range internalNameServers {
		nameServers = append(nameServers, NameServer{Host: ns.Host})
	}
	return nameServers, nil
}

// NetworkByHost returns the network address for the provided hostname.
func NetworkByHost(host string) (*net.IP, error) {
	ipAddr := net.ParseIP(host)
	if ipAddr == nil {
		addr, err := AddrByHostName(host)
		if err != nil {
			return nil, xerrors.Errorf("%q is an invalild host: %v", host, err)
		}
		ipAddr = *addr
	}

	network := ipAddr.Mask(ipAddr.DefaultMask())
	if network == nil {
		return nil, xerrors.Errorf("failed to get network address of host %q", ipAddr.String())
	}
	return &network, nil
}

// HostAndAddr returns the hostname and ip address of host whether host is an IP address or a hostname.
// HostAndAddr returns a non-nil error if host is an invalid ip address or a hostname that cannot be resolved to an IP address.
func HostAndAddr(host string) (string, *net.IP, error) {
	var hostname string

	ip := net.ParseIP(host)
	if ip == nil {
		hostname = host
		addr, err := AddrByHostName(hostname)
		if err != nil {
			return "", nil, xerrors.Errorf("failed to get ip address by hostname %q: %w", hostname, err)
		}
		ip = *addr
	} else {
		host, err := HostNameByIP(ip)
		if err != nil {
			return "", nil, xerrors.Errorf("failed to get hostname by ip address %q: %w", ip, err)
		}
		hostname = host
	}
	return hostname, &ip, nil
}

// ServiceProvider returns the internet service provider information for ip.
func ServiceProvider(ip *net.IP) (*InternetServiceProvider, error) {
	client, err := ipisp.NewDNSClient()
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize new dns client: %w", err)
	}
	defer client.Close()

	resp, err := client.LookupIP(*ip)
	if err != nil {
		return nil, err
	}

	return &InternetServiceProvider{
		Name:                    resp.Name.Raw,
		IP:                      &resp.IP,
		Country:                 resp.Country,
		Registry:                resp.Registry,
		IpRange:                 resp.Range,
		AutonomousServiceNumber: resp.ASN.String(),
		AllocatedAt:             &resp.AllocatedAt,
	}, nil
}
