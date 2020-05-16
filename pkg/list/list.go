package list

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	gw "github.com/jackpal/gateway"
)

const (
	googleDNS = "8.8.8.8:80"
	getExtURL = "http://myexternalip.com/raw"
	notFound  = "not found"
)

var stars = strings.Repeat("*", 30)

type pong struct {
	name string
	ip   string
	up   bool
}

// Me prints out the device name, remote, and local IP addresses of this machine.
// It also prints out the router IP.
func Me() error {
	if err := localIP(); err != nil {
		return err
	}

	if err := remoteIP(); err != nil {
		return err
	}

	if err := router(); err != nil {
		return err
	}
	return nil
}

func localIP() error {
	conn, err := net.Dial("udp", googleDNS)
	if err != nil {
		return fmt.Errorf("failed to dial google dns : %s", err)
	}
	defer conn.Close()
	localAddr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return fmt.Errorf("failed to resolve local IP")
	}
	host, _, err := net.SplitHostPort(localAddr.String())
	if err != nil {
		return err
	}
	names, _ := net.LookupAddr(host)
	if len(names) > 0 {
		fmt.Printf("Name : %s\n", names[0])
	}
	fmt.Printf("Local IP : %s\n", host)
	return nil
}

func remoteIP() error {
	resp, err := http.Get(getExtURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Remote IP : %s\n", string(data))
	return nil
}

func router() error {
	gatewayAddr, err := gw.DiscoverGateway()
	if err != nil {
		return err
	}
	fmt.Printf("Gateway : %s\n", gatewayAddr.String())
	return nil
}

// AllDevices lists IP address, name, and host of all connected network devices.
func AllDevices() error {
	hosts, err := hosts("192.168.1.0/24")
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(hosts))
	for _, h := range hosts {
		go func(host string) {
			process(host)
			wg.Done()
		}(h)
	}
	wg.Wait()
	return nil
}

func process(host string) {
	var up bool

	_, err := net.DialTimeout("ip", host, time.Second)
	if err != nil {
		up = false
	} else {
		up = true
	}

	names, err := net.LookupAddr(host)
	if err != nil {
		return
	}
	if len(names) > 0 {
		fmt.Println(stars)
		fmt.Printf("Host : %s\nIP : %s\nConnected : %t\n", names[0], host, up)
	}
}

func hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func match(a, b string) bool {
	return a == b
}
