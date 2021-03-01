package list

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"

	"github.com/fuskovic/networker/internal/lookup"
	gw "github.com/jackpal/gateway"
	"golang.org/x/xerrors"
)

// Devices lists all of the devices on the local network.
func Devices(ctx context.Context) ([]*Device, error) {
	cidr, err := getCIDR(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to get cidr: %w", err)
	}

	router, err := getRouter(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to get router ip: %w", err)
	}

	hostIPs, err := getHosts(ctx, cidr, router)
	if err != nil {
		return nil, xerrors.Errorf("failed to get hosts: %w", err)
	}

	currentDevice, err := getCurrentDevice(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to get current device: %w", err)
	}

	var (
		devices = []*Device{currentDevice, router}
		wg      = sync.WaitGroup{}
		mutex   = sync.Mutex{}
	)

	for _, hostIP := range hostIPs {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			device, err := getDevice(ctx, ip)
			if err != nil || device == nil {
				return
			}
			mutex.Lock()
			devices = append(devices, device)
			mutex.Unlock()
		}(hostIP)
	}
	wg.Wait()
	return devices, nil
}

func getDevice(_ context.Context, ip string) (*Device, error) {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return nil, xerrors.Errorf("failed to parse ip %q", ip)
	}

	hostname, err := lookup.HostNameByIP(ipAddr)
	if err != nil {
		return nil, xerrors.Errorf("failed to lookup hostname by ip address %q: %w", ip, err)
	}
	return &Device{
		localIP:  ipAddr,
		hostname: hostname,
		kind:     DeviceKindPeer,
	}, nil
}

func getCurrentDevice(_ context.Context) (*Device, error) {
	localIP, err := getLocalIP()
	if err != nil {
		return nil, xerrors.Errorf("failed to get local ip of current device: %w", err)
	}

	remoteIP, err := getRemoteIP()
	if err != nil {
		return nil, xerrors.Errorf("failed to get remote ip of current device: %w", err)
	}

	hostname, err := lookup.HostNameByIP(localIP)
	if err != nil {
		return nil, xerrors.Errorf("failed to get host: %w", err)
	}
	return &Device{
		localIP:  localIP,
		remoteIP: remoteIP,
		hostname: hostname,
		kind:     DeviceKindCurrent,
	}, nil
}

func getRouter(_ context.Context) (*Device, error) {
	ipAddr, err := gw.DiscoverGateway()
	if err != nil {
		return nil, xerrors.Errorf("failed to discover gateway: %w", err)
	}
	hostname, err := lookup.HostNameByIP(ipAddr)
	if err != nil {
		return nil, xerrors.Errorf("failed to get router name: %w", err)
	}
	return &Device{
		hostname: hostname,
		localIP:  ipAddr,
		kind:     DeviceKindRouter,
	}, nil
}

func getCIDR(_ context.Context) (string, error) {
	localIP, err := getLocalIP()
	if err != nil {
		return "", xerrors.Errorf("failed to get local ip: %w", err)
	}
	return fmt.Sprintf("%s/24", localIP.Mask(localIP.DefaultMask())), nil
}

func getHosts(_ context.Context, cidr string, router *Device) ([]string, error) {
	ip, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse cidr %s: %w", cidr, err)
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
		if ip.String() == router.localIP.String() {
			continue
		}
		ips = append(ips, ip.String())
	}
	return ips[1 : len(ips)-1], nil
}

func getLocalIP() (net.IP, error) {
	c, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, xerrors.Errorf("failed to dial google dns : %w", err)
	}
	defer c.Close()

	localAddr, ok := c.LocalAddr().(*net.UDPAddr)
	if !ok {
		return nil, xerrors.New("failed to resolve local IP")
	}
	return localAddr.IP, nil
}

func getRemoteIP() (net.IP, error) {
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
		return nil, xerrors.Errorf("failed to get resolve ip %q as ipv4", b)
	}
	return remoteIP, nil
}
