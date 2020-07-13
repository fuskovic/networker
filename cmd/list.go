package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	gw "github.com/jackpal/gateway"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"
)

const (
	googleDNS = "8.8.8.8:80"
	getExtURL = "http://myexternalip.com/raw"
)

type (
	status struct {
		addr, name string
		connected  bool
	}
	listCmd struct{ me, all bool }
)

// Spec returns a command spec containing a description of it's usage.
func (cmd *listCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "list",
		Usage: "[flags]",
		Desc:  "List information on connected network devices.",
	}
}

// RegisterFlags initializes how a flag set is processed for a particular command.
func (cmd *listCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.BoolVarP(&cmd.me, "me", "m", cmd.me, "Lists the local and remote IP of this machine and the router IP.")
	fl.BoolVarP(&cmd.all, "all", "a", cmd.all, "List the IP, hostname, and connection status of all devices on this network.")
}

// Run prints either general network information for this machine or for the entire network
// depending on how the flag set has been configured.
func (cmd *listCmd) Run(fl *pflag.FlagSet) {
	ctx := context.Background()

	switch {
	case cmd.me:
		me(ctx)
	case cmd.all:
		if err := all(ctx); err != nil {
			flog.Error("failed to list all network devices : %v", err)
			fl.Usage()
		}
	default:
		fl.Usage()
	}
}

func (s *status) fields() []slog.Field {
	return []slog.Field{
		slog.F("addr", s.addr),
		slog.F("name", s.name),
		slog.F("connected", s.connected),
	}
}

func me(ctx context.Context) {
	log := sloghuman.Make(os.Stdout)
	var fields []slog.Field

	local, err := localIP()
	if err != nil {
		local = "unknown"
	}
	fields = append(fields, slog.F("local", local))

	remote, err := remoteIP()
	if err != nil {
		remote = "unknown"
	}
	fields = append(fields, slog.F("remote", remote))

	router, err := router()
	if err != nil {
		router = "unknown"
	}
	fields = append(fields, slog.F("router", router))

	log.Info(ctx, "me", fields...)
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

func all(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	cidr, err := getCIDR()
	if err != nil {
		return err
	}

	hosts, err := hosts(cidr)
	if err != nil {
		return err
	}

	log := sloghuman.Make(os.Stdout)
	statusChan := make(chan status)

	for _, h := range hosts {
		go func(host string) {
			process(ctx, host, statusChan)
		}(h)
	}

	var numDevices int

	for {
		select {
		case <-ctx.Done():
			log.Info(ctx, "complete", slog.F("devices", numDevices))
			return nil
		case s := <-statusChan:
			numDevices++
			log.Info(ctx, "device", s.fields()...)
		}
	}
}

func process(ctx context.Context, host string, sc chan<- status) {
	var up bool

	if _, err := net.DialTimeout("ip", host, 100*time.Millisecond); err == nil {
		up = true
	}

	names, err := net.LookupAddr(host)
	if err != nil {
		return
	}

	if len(names) > 0 {
		sc <- status{
			addr:      host,
			name:      names[0],
			connected: up,
		}
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
