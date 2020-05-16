package list

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"strings"

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hosts, err := hosts("192.168.1.0/24")
	if err != nil {
		return err
	}
	concurrentMax := 100
	pingChan := make(chan string, concurrentMax)
	pongChan := make(chan pong, len(hosts))
	doneChan := make(chan pong)

	for i := 0; i < concurrentMax; i++ {
		go ping(pingChan, pongChan)
	}

	go receivePong(cancel, len(hosts), pongChan, doneChan)

	for _, ip := range hosts {
		pingChan <- ip
	}

processing:
	for {
		select {
		case d := <-doneChan:
			fmt.Println(stars)
			fmt.Printf("Name: %s\nIP: %s\nconnected: %t\n", d.name, d.ip, d.up)
		case <-ctx.Done():
			break processing
		default:
			continue
		}
	}
	return nil
}

func ping(pingChan <-chan string, pongChan chan<- pong) {
	var alive bool
	var host string

	for ip := range pingChan {
		if _, err := exec.Command("ping", "-c1", "-t1", ip).Output(); err != nil {
			alive = false
		} else {
			alive = true
		}
		names, _ := net.LookupAddr(ip)
		if len(names) > 0 {
			host = names[0]
		} else {
			host = notFound
		}
		pongChan <- pong{host, ip, alive}
	}
}

func receivePong(cancel context.CancelFunc, pongNum int, pongChan <-chan pong, doneChan chan<- pong) {
	for i := 0; i < pongNum; i++ {
		pong := <-pongChan
		if pong.name != notFound {
			doneChan <- pong
		}
	}
	cancel()
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
