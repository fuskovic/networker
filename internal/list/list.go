package list

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	gw "github.com/jackpal/gateway"
	goping "github.com/tatsushid/go-fastping"

	"github.com/fuskovic/networker/v3/internal/resolve"
)

const (
	notAvailable = "N/A"
	icmpIpv4     = "ip4:icmp"
	icmpIpv6     = "ip6:ipv6-icmp"
)

const (
	DeviceKindUnknown Kind = "unknown"
	DeviceKindRouter  Kind = "router"
	DeviceKindCurrent Kind = "current-device"
	DeviceKindPeer    Kind = "peer"
)

type Kind string

type Device struct {
	Kind     Kind   `json:"kind" table:"KIND"`
	Hostname string `json:"hostname" table:"HOSTNAME"`
	LocalIP  net.IP `json:"local_ip" table:"LOCAL_IP"`
	RemoteIP net.IP `json:"remote_ip,omitempty" table:"REMOTE_IP"`
	Up       bool   `json:"up" yaml:"up" table:"UP"`
}

// Devices lists all of the devices on the local network.
func Devices(ctx context.Context) ([]Device, error) {
	cidr, err := getCIDR(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get cidr: %w", err)
	}

	router, err := getRouter(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get router: %w", err)
	}

	hostIPs, err := getHosts(ctx, cidr, router)
	if err != nil {
		return nil, fmt.Errorf("failed to get hosts: %w", err)
	}

	currentDevice, err := getCurrentDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current device: %w", err)
	}

	hostIPs = removeIP(currentDevice.LocalIP.String(), hostIPs)

	var (
		devices = []Device{*router, *currentDevice}
		wg      = sync.WaitGroup{}
		mutex   = sync.Mutex{}
	)

	for _, hostIP := range dedupe(hostIPs) {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()

			device, err := getDevice(ctx, ip)
			if err != nil || device == nil {
				return
			}
			device.Up = ping(ip)

			mutex.Lock()
			devices = append(devices, *device)
			mutex.Unlock()
		}(hostIP)
	}
	wg.Wait()
	return devices, nil
}

func getDevice(_ context.Context, ip string) (*Device, error) {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return nil, fmt.Errorf("failed to parse ip %q", ip)
	}

	return &Device{
		LocalIP:  ipAddr,
		Hostname: resolve.Hostname(ipAddr),
		Kind:     DeviceKindPeer,
	}, nil
}

func getCurrentDevice(_ context.Context) (*Device, error) {
	localIP, err := getCurrentDeviceLocalIP()
	if err != nil {
		return nil, fmt.Errorf("failed to get local ip of current device: %w", err)
	}

	remoteIP, err := getCurrentDeviceRemoteIP()
	if err != nil {
		return nil, fmt.Errorf("failed to get remote ip of current device: %w", err)
	}

	return &Device{
		LocalIP:  localIP,
		RemoteIP: remoteIP,
		Hostname: resolve.Hostname(localIP),
		Kind:     DeviceKindCurrent,
		Up:       true,
	}, nil
}

func getRouter(_ context.Context) (*Device, error) {
	ipAddr, err := gw.DiscoverGateway()
	if err != nil {
		return nil, fmt.Errorf("failed to discover gateway: %w", err)
	}

	return &Device{
		Hostname: resolve.Hostname(ipAddr),
		LocalIP:  ipAddr,
		Kind:     DeviceKindRouter,
		Up:       true,
	}, nil
}

func getCIDR(_ context.Context) (string, error) {
	localIP, err := getCurrentDeviceLocalIP()
	if err != nil {
		return "", fmt.Errorf("failed to get local ip: %w", err)
	}
	return fmt.Sprintf("%s/24", localIP.Mask(localIP.DefaultMask())), nil
}

func getHosts(_ context.Context, cidr string, router *Device) ([]string, error) {
	ip, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cidr %s: %w", cidr, err)
	}

	inc := func(ip net.IP) {
		for j := len(ip) - 1; j >= 0; j-- {
			ip[j]++
			if ip[j] > 0 {
				break
			}
		}
	}

	var ips []string
	for ip := ip.Mask(network.Mask); network.Contains(ip); inc(ip) {
		if ip.String() == router.LocalIP.String() {
			continue
		}
		ips = append(ips, ip.String())
	}
	return ips[1 : len(ips)-1], nil
}

func getCurrentDeviceLocalIP() (net.IP, error) {
	c, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, fmt.Errorf("failed to dial google dns : %w", err)
	}
	defer c.Close()

	localAddr, ok := c.LocalAddr().(*net.UDPAddr)
	if !ok {
		return nil, errors.New("failed to resolve local IP")
	}
	return localAddr.IP, nil
}

func getCurrentDeviceRemoteIP() (net.IP, error) {
	r, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var remoteIP net.IP
	ipAddr := string(b)
	if strings.Contains(ipAddr, ":") {
		remoteIP = net.ParseIP(ipAddr).To16()
	} else {
		remoteIP = net.ParseIP(ipAddr).To4()
	}
	if remoteIP == nil {
		return nil, fmt.Errorf("failed to get resolve ip %q as ipv4", b)
	}
	return remoteIP, nil
}

func dedupe(hosts []string) []string {
	m := make(map[string]int)
	var filteredHosts []string
	for _, host := range hosts {
		if m[host] == 0 {
			m[host]++
			filteredHosts = append(filteredHosts, host)
			continue
		}
	}
	return filteredHosts
}

func removeIP(ip string, hosts []string) []string {
	var filteredHosts []string
	for _, host := range hosts {
		if host == ip {
			continue
		}
		filteredHosts = append(filteredHosts, host)
	}
	return filteredHosts
}

func ping(ip string) bool {
	p := goping.NewPinger()
	p.Network("udp")

	netProto := icmpIpv4
	if strings.Contains(ip, ":") {
		netProto = icmpIpv6
	}

	addr, err := net.ResolveIPAddr(netProto, ip)
	if err != nil {
		return false
	}

	p.AddIPAddr(addr)
	p.MaxRTT = time.Second

	var up bool
	p.OnRecv = func(addr *net.IPAddr, t time.Duration) {
		up = true
	}
	p.Run()
	return up
}
