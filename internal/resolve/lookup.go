package resolve

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/ammario/ipisp"
)

// Record can be used as a common type between lookup commands
// that supports json and table output.
type Record struct {
	Hostname string `json:"hostname" yaml:"hostname" table:"HOSTNAME"`
	IP       net.IP `json:"ip" yaml:"ip" table:"IP_ADDRESS"`
}

type NetworkRecord struct {
	Hostname  string `json:"hostname" yaml:"hostname" table:"HOSTNAME"`
	NetworkIP net.IP `json:"network" yaml:"network" table:"NETWORK"`
}

type Func func(string) error

// InternetServiceProvider describes an internet service provider.
type InternetServiceProvider struct {
	Name                    string     `json:"name" table:"NAME"`
	IP                      *net.IP    `json:"ip_address" table:"IP"`
	Country                 string     `json:"country" table:"COUNTRY"`
	Registry                string     `json:"registry" table:"REGISTRY"`
	IpRange                 *net.IPNet `json:"ip_range" table:"IP_RANGE"`
	AutonomousServiceNumber string     `json:"autonomous_service_number" table:"ASN"`
	AllocatedAt             *time.Time `json:"allocated_at" table:"ALLOCATED_AT"`
}

// NameServer is used in place of the standard library object to support table writes.
type NameServer struct {
	IP   net.IP `json:"ip" table:"IP"`
	Host string `json:"nameserver" table:"Nameserver"`
}

// HostNameByIP returns the hostname for the provided ip address.
func HostNameByIP(ip net.IP) (*Record, error) {
	hostnames, err := HostNamesByIP(ip)
	if err != nil {
		return nil, err
	}
	return &Record{
		IP:       ip,
		Hostname: hostnames[0],
	}, nil
}

// HostNamesByIP returns all hostnames found for the provided ip address.
func HostNamesByIP(ip net.IP) ([]string, error) {
	hostnames, err := net.LookupAddr(ip.String())
	if err != nil {
		return nil, fmt.Errorf("failed to lookup hostnames for ip address %q : %v", ip, err)
	}
	if len(hostnames) == 0 {
		return nil, fmt.Errorf("no hostnames found for ip address: %q", ip)
	}
	return hostnames, nil
}

// AddrByHostName resolves the ip address of the provided hostname.
func AddrByHostName(hostname string) (*Record, error) {
	ipAddrs, err := AddrsByHostName(hostname)
	if err != nil {
		return nil, err
	}

	if len(ipAddrs) == 0 {
		return nil, fmt.Errorf("no addresses found for hostname %q", hostname)
	}

	ipv4 := ipAddrs[0].To4()
	if ipv4 == nil {
		return &Record{
			Hostname: hostname,
			IP:       *ipAddrs[0],
		}, nil
	}

	return &Record{
		Hostname: hostname,
		IP:       ipv4,
	}, err
}

// AddrsByHostName returns all ip addresses found for the provided hostname.
func AddrsByHostName(hostname string) ([]*net.IP, error) {
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		return nil, fmt.Errorf("failed to look up ip addresses for hostname %q: %v", hostname, err)
	}

	if len(addrs) == 0 {
		return nil, fmt.Errorf("no ip addresses found for hostname %q", hostname)
	}

	var ipAddrs []*net.IP

	for _, a := range addrs {
		ipAddr := net.ParseIP(a).To4()
		if ipAddr == nil {
			continue
		}
		ipAddrs = append(ipAddrs, &ipAddr)
	}
	return ipAddrs, nil
}

// NameServersByHostName looks up all nameservers for the provided hostname.
func NameServersByHostName(hostname string) ([]NameServer, error) {
	internalNameServers, err := net.LookupNS(stripHostname(hostname))
	if err != nil {
		return nil, fmt.Errorf("failed to look up name server for hostname %q : %v", hostname, err)
	}
	if len(internalNameServers) == 0 {
		return nil, fmt.Errorf("no name servers found for hostname %q", hostname)
	}
	var nameServers []NameServer
	for _, ns := range internalNameServers {
		record, err := AddrByHostName(ns.Host)
		if err != nil {
			return nil, fmt.Errorf("failed to get ip by hostname %q: %w", ns.Host, err)
		}
		nameServers = append(nameServers,
			NameServer{
				IP:   record.IP.To4(),
				Host: ns.Host,
			},
		)
	}
	return nameServers, nil
}

// NetworkByHost returns the network address for the provided hostname.
func NetworkByHost(host string) (*NetworkRecord, error) {
	var hostname string
	ipAddr := net.ParseIP(host)
	if ipAddr == nil {
		hostname = host
		record, err := AddrByHostName(host)
		if err != nil {
			return nil, fmt.Errorf("%q is an invalild host: %v", host, err)
		}
		ipAddr = record.IP
	} else {
		record, err := HostNameByIP(ipAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve hostname for ip %q: %s", ipAddr, err)
		}
		hostname = record.Hostname
	}

	network := ipAddr.Mask(ipAddr.DefaultMask())
	if network == nil {
		return nil, fmt.Errorf("failed to get network address of host %q", ipAddr.String())
	}

	return &NetworkRecord{
		Hostname:  hostname,
		NetworkIP: network,
	}, nil
}

// HostAndAddr returns the hostname and ip address of host whether host is an IP address or a hostname.
// HostAndAddr returns a non-nil error if host is an invalid ip address or a hostname that cannot be resolved to an IP address.
func HostAndAddr(host string) (string, *net.IP, error) {
	var hostname string

	ip := net.ParseIP(host)
	if ip == nil {
		hostname = host
		record, err := AddrByHostName(hostname)
		if err != nil {
			return "", nil, fmt.Errorf("failed to get ip address by hostname %q: %w", hostname, err)
		}
		ip = record.IP
	} else {
		record, err := HostNameByIP(ip)
		if err != nil {
			return "", nil, fmt.Errorf("failed to get hostname by ip address %q: %w", ip, err)
		}
		hostname = record.Hostname
	}
	return hostname, &ip, nil
}

// ServiceProvider returns the internet service provider information for ip.
func ServiceProvider(ip *net.IP) (*InternetServiceProvider, error) {
	client, err := ipisp.NewDNSClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize new dns client: %w", err)
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

func stripHostname(hostname string) string {
	hostname = strings.ReplaceAll(hostname, "https://", "")
	hostname = strings.ReplaceAll(hostname, "http://", "")
	hostname = strings.ReplaceAll(hostname, "www.", "")
	return hostname
}
