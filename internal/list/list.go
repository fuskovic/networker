package list

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	u "github.com/fuskovic/networker/internal/utils"
	gw "github.com/jackpal/gateway"
)

const (
	local addrType = iota
	remote
	router
)

type (
	addrFunc func() (string, error)
	addrType int
)

// String returns the string description of the addrType.
//
// 0 - local
// 1 - remote
// 2 - router
// default : unknown
func (a *addrType) String() string {
	var s string
	switch *a {
	case local:
		s = "local"
	case remote:
		s = "remote"
	case router:
		s = "router"
	default:
		s = "unknown"
	}
	return s
}

// List lists the network information of all the devices on the current network.
func List(ctx context.Context) error {
	t := time.Duration(3) * time.Second
	ctx, cancel := context.WithTimeout(ctx, t)
	defer cancel()

	log := sloghuman.Make(os.Stdout)
	log.Info(ctx, "current device", current(ctx)...)

	cidr, err := getCIDR()
	if err != nil {
		return err
	}

	hosts, err := hosts(cidr)
	if err != nil {
		return err
	}

	rc := make(chan u.Row)

	for _, h := range hosts {
		go func(h string) {
			process(ctx, h, rc)
		}(h)
	}

	var numDevices int
	log.Info(ctx, "LAN")

	for {
		select {
		case <-ctx.Done():
			log.Info(ctx, "found", slog.F(
				"devices", numDevices),
			)
			return nil
		case r := <-rc:
			numDevices++
			log.Info(ctx, "device", r...)
		}
	}
}

func current(ctx context.Context) u.Row {
	var r u.Row

	host := func(addr string) string {
		hostnames, err := net.LookupAddr(addr)
		if err != nil {
			return ""
		}

		if len(hostnames) < 1 {
			return ""
		}
		return hostnames[0]
	}

	funcs := []addrFunc{
		localIP,
		remoteIP,
		routerIP,
	}

	for i, f := range funcs {
		a, err := f()
		if err != nil {
			a = "unknown"
		}
		at := addrType(i)
		r.Add(at.String(), a)
		if at == local {
			r.Add("hostname", host(a))
		}
	}
	return r
}

func localIP() (string, error) {
	c, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", fmt.Errorf("failed to dial google dns : %s", err)
	}
	defer c.Close()

	a, ok := c.LocalAddr().(*net.UDPAddr)
	if !ok {
		return "", fmt.Errorf("failed to resolve local IP")
	}

	h, _, err := net.SplitHostPort(a.String())
	if err != nil {
		return "", err
	}
	return h, nil
}

func remoteIP() (string, error) {
	r, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func routerIP() (string, error) {
	a, err := gw.DiscoverGateway()
	if err != nil {
		return "", err
	}
	return a.String(), nil
}

func process(ctx context.Context, h string, rc chan<- u.Row) {
	names, err := net.LookupAddr(h)
	if err != nil {
		return
	}

	if len(names) > 0 {
		var r u.Row
		r.Add("addr", h)
		r.Add("hostname", names[0])
		rc <- r
	}
}

func hosts(cidr string) ([]string, error) {
	ip, n, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(n.Mask); n.Contains(ip); inc(ip) {
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

func getCIDR() (string, error) {
	host, err := localIP()
	if err != nil {
		return "", err
	}

	a := net.ParseIP(host)
	if a == nil {
		return "", err
	}
	m := a.DefaultMask()
	return fmt.Sprintf("%s/24", a.Mask(m)), nil
}
