package list

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"

	"github.com/fuskovic/networker/internal/resolve"
	gw "github.com/jackpal/gateway"
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

	record, err := resolve.HostNameByIP(ipAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup hostname by ip address %q: %w", ip, err)
	}

	return &Device{
		LocalIP:  ipAddr,
		Hostname: record.Hostname,
		Kind:     DeviceKindPeer,
	}, nil
}

func getCurrentDevice(_ context.Context) (*Device, error) {
	localIP, err := getCurrentDeviceLocalIP()
	if err != nil {
		return nil, fmt.Errorf("failed to get local ip of current device: %w", err)
	}

	remoteIP, err := getCurrentDeviceRemoteIP(localIP)
	if err != nil {
		return nil, fmt.Errorf("failed to get remote ip of current device: %w", err)
	}

	record, err := resolve.HostNameByIP(localIP)
	if err != nil {
		return nil, fmt.Errorf("failed to get host: %w", err)
	}

	return &Device{
		LocalIP:  localIP,
		RemoteIP: remoteIP,
		Hostname: record.Hostname,
		Kind:     DeviceKindCurrent,
	}, nil
}

func getRouter(_ context.Context) (*Device, error) {
	ipAddr, err := gw.DiscoverGateway()
	if err != nil {
		return nil, fmt.Errorf("failed to discover gateway: %w", err)
	}

	record, err := resolve.HostNameByIP(ipAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve hostname by ip for gateway: %w", err)
	}

	return &Device{
		Hostname: record.Hostname,
		LocalIP:  ipAddr,
		Kind:     DeviceKindRouter,
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

func getCurrentDeviceRemoteIP(localIP net.IP) (net.IP, error) {
	r, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	remoteIP := net.ParseIP(string(b)).To4()
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
