package cmd

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os/user"
	"strings"
	"sync"
	"time"

	gw "github.com/jackpal/gateway"
	fp "github.com/tatsushid/go-fastping"

	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"
)

const (
	googleDNS = "8.8.8.8:80"
	getExtURL = "http://myexternalip.com/raw"
)

type listCmd struct{ me, all bool }

// Spec returns a command spec containing a description of it's usage.
func (cmd *listCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "list",
		Usage: "TODO: ADD USAGE",
		Desc:  "List information on connected network devices.",
	}
}

// RegisterFlags initializes how a flag set is processed for a particular command.
func (cmd *listCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.BoolVarP(&cmd.me, "me", "m", cmd.me, "List the name, local IP, remote IP, and router IP for this machine and the network it's connected to.")
	fl.BoolVarP(&cmd.all, "all", "a", cmd.all, "List the IP, hostname, and connection status of all devices on this network. (must be run as root)")
}

// Run prints either general network information for this machine or for the entire network
// depending on how the flag set has been configured.
func (cmd *listCmd) Run(fl *pflag.FlagSet) {
	u, err := user.Current()
	if err != nil {
		flog.Fatal(err.Error())
	}
	switch {
	case u.Uid != "0":
		flog.Fatal("Must be root to run this command")
	case cmd.me:
		me()
	case cmd.all:
		all()
	default:
		fl.Usage()
	}
}

func me() {
	local, err := localIP()
	if err != nil {
		local = "unknown"
	}

	remote, err := remoteIP()
	if err != nil {
		remote = "unknown"
	}

	router, err := router()
	if err != nil {
		router = "unknown"
	}

	sep := strings.Repeat(" ", 11)
	flog.Info("LOCAL%sREMOTE%sROUTER", sep, sep)
	flog.Info("%s   %s   %s", local, remote, router)
}

func localIP() (string, error) {
	conn, err := net.Dial("udp", googleDNS)
	if err != nil {
		return "", fmt.Errorf("failed to dial google dns : %s", err)
	}
	defer conn.Close()

	localAddr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return "", fmt.Errorf("failed to resolve local IP")
	}

	host, _, err := net.SplitHostPort(localAddr.String())
	if err != nil {
		return "", err
	}
	return host, nil
}

func remoteIP() (string, error) {
	resp, err := http.Get(getExtURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func router() (string, error) {
	gatewayAddr, err := gw.DiscoverGateway()
	if err != nil {
		return "", err
	}
	return gatewayAddr.String(), nil
}

func all() error {
	cidr, err := getCIDR()
	if err != nil {
		return err
	}

	hosts, err := hosts(cidr)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(hosts))

	sep := func(n int) string { return strings.Repeat(" ", n) }
	flog.Info("IP%sHOST%sCONNECTED", sep(12), sep(37))

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

	p := fp.NewPinger()
	ip, err := net.ResolveIPAddr("ip4:icmp", host)
	if err != nil {
		flog.Error(err.Error())
		return
	}
	p.AddIPAddr(ip)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		up = true
	}
	p.Run()

	names, err := net.LookupAddr(host)
	if err != nil {
		return
	}

	if len(names) > 0 {
		for len(names[0]) < 40 {
			names[0] += " "
		}
		flog.Info("%s %s %t", host, names[0], up)
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

func getCIDR() (string, error) {
	host, err := localIP()
	if err != nil {
		return "", err
	}

	addr := net.ParseIP(host)
	if addr == nil {
		return "", err
	}

	mask := addr.DefaultMask()
	cidr := fmt.Sprintf("%s/24", addr.Mask(mask))
	return cidr, nil
}
