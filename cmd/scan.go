package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"
)

const (
	// TotalPorts is the total number of all tcp/udp ports
	TotalPorts = 65535
	udp        = "udp"
)

var timeOut = 3 * time.Second

type (
	scanCmd struct {
		addr                       string
		ports                      []int
		upTo                       int
		tCPonly, uDPonly, openOnly bool
	}
	result struct {
		protocol, addr string
		port           int
		open           bool
	}
)

// Spec returns a command spec containing a description of it's usage.
func (cmd *scanCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "scan",
		Usage: "[flags]",
		Desc:  "Scan a host for exposed ports.",
	}
}

func (cmd *scanCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVar(&cmd.addr, "ip", "", "IP address to scan.")
	fl.IntSliceVarP(&cmd.ports, "ports", "p", cmd.ports, "Specify a comma-separated list of ports to scan. (scans all ports if left unspecified)")
	fl.IntVarP(&cmd.upTo, "up-to", "u", cmd.upTo, "Scan all ports up to a given port number.")
	fl.BoolVarP(&cmd.tCPonly, "tcp-only", "t", cmd.tCPonly, "Only scan TCP ports.")
	fl.BoolVar(&cmd.uDPonly, "udp-only", cmd.uDPonly, "Only scan UDP ports.")
	fl.BoolVarP(&cmd.openOnly, "open-only", "o", cmd.openOnly, "Only print the ports that are open.")
}

func (cmd *scanCmd) Run(fl *pflag.FlagSet) {
	ctx := context.Background()

	if net.ParseIP(cmd.addr) == nil {
		flog.Error("%s is not a valid IP address", cmd.addr)
		fl.Usage()
		return
	}

	if !protocolSpecified(cmd.tCPonly, cmd.uDPonly) {
		flog.Info("protocol unspecified enabling scanner for both")
		cmd.tCPonly = true
		cmd.uDPonly = true
	}

	switch {
	case cmd.upTo > TotalPorts:
		flog.Error("can not scan more than %d ports", TotalPorts)
		fl.Usage()
		return
	case len(cmd.ports) > 0:
		cmd.scanPorts(ctx, cmd.ports)
	case cmd.upTo > 0:
		cmd.scanUpTo(ctx, cmd.upTo)
	default:
		cmd.scanAllPorts(ctx)
	}
	flog.Success("scan complete")
}

func (cmd *scanCmd) scanPorts(ctx context.Context, specifiedPorts []int) {
	cmd.start(ctx, specifiedPorts)
}

func (cmd *scanCmd) scanUpTo(ctx context.Context, upTo int) {
	portsForScanning := portsToScan(upTo)
	cmd.start(ctx, portsForScanning)
}

func (cmd *scanCmd) scanAllPorts(ctx context.Context) {
	portsForScanning := portsToScan(TotalPorts)
	cmd.start(ctx, portsForScanning)
}

func (cmd *scanCmd) scan(port int, c chan<- result) {
	send := func(protocol string) {
		c <- result{
			addr:     cmd.addr,
			port:     port,
			protocol: protocol,
			open:     isOpen(protocol, cmd.addr, port),
		}
	}

	if cmd.tCPonly {
		if cmd.shouldLog(tcp, port) {
			send(tcp)
		}
	}

	if cmd.uDPonly {
		if cmd.shouldLog(udp, port) {
			send(udp)
		}
	}
}

func (cmd *scanCmd) shouldLog(protocol string, port int) bool {
	return cmd.openOnly && isOpen(protocol, cmd.addr, port) || !cmd.openOnly
}

func (r *result) fields() []slog.Field {
	return []slog.Field{
		slog.F("protocol", r.protocol),
		slog.F("port", r.port),
		slog.F("open", r.open),
	}
}

func isOpen(protocol, host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout(protocol, address, timeOut)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func portsToScan(max int) (ports []int) {
	for port := 0; port < max; port++ {
		ports = append(ports, port)
	}
	return
}

func protocolSpecified(tcp, udp bool) bool {
	return tcp == true || udp == true
}

func organize(results []result) []result {
	var organized []result
	for i := 0; i < len(results); i++ {
		for _, r := range results {
			if r.port == i {
				organized = append(organized, r)
			}
		}
	}
	return organized
}

func (cmd *scanCmd) start(ctx context.Context, portsForScanning []int) {
	var (
		log     = sloghuman.Make(os.Stdout)
		wg      sync.WaitGroup
		results []result
	)

	wg.Add(len(portsForScanning))
	ch := make(chan result)

	go func() {
		for {
			select {
			case r, ok := <-ch:
				results = append(results, r)
				if !ok {
					return
				}
			}
		}
	}()

	for _, port := range portsForScanning {
		go func(p int) {
			cmd.scan(p, ch)
			wg.Done()
		}(port)
	}
	wg.Wait()
	close(ch)

	for _, r := range organize(results) {
		log.Info(ctx, r.addr, r.fields()...)
	}
}
