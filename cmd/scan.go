package cmd

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"
)

const (
	// TotalPorts is the total number of all tcp/udp ports
	TotalPorts = 65535
	udp        = "udp"
	stars      = ""
)

var timeOut = 3 * time.Second

type scanCmd struct {
	addr                       string
	ports                      []int
	upTo                       int
	tCPonly, uDPonly, openOnly bool
}

// Spec returns a command spec containing a description of it's usage.
func (cmd *scanCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "scan",
		Usage: "TODO: ADD USAGE",
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
	if net.ParseIP(cmd.addr) == nil {
		flog.Fatal(fmt.Sprintf("%s is not a valid IP address", cmd.addr))
	}

	if !protocolSpecified(cmd.tCPonly, cmd.uDPonly) {
		flog.Info("protocol unspecified enabling scanner for both")
		cmd.tCPonly = true
		cmd.uDPonly = true
	}

	switch {
	case cmd.upTo > TotalPorts:
		flog.Fatal(fmt.Sprintf("can not scan more than %d ports", TotalPorts))
	case len(cmd.ports) > 0:
		cmd.scanPorts(cmd.ports)
	case cmd.upTo > 0:
		cmd.scanUpTo(cmd.upTo)
	default:
		cmd.scanAllPorts()
	}
	flog.Success("scan complete")
}

func (cmd *scanCmd) scanPorts(specifiedPorts []int) {
	cmd.start(specifiedPorts)
}

func (cmd *scanCmd) scanUpTo(upTo int) {
	portsForScanning := portsToScan(upTo)
	cmd.start(portsForScanning)
}

func (cmd *scanCmd) scanAllPorts() {
	portsForScanning := portsToScan(TotalPorts)
	cmd.start(portsForScanning)
}

func (cmd *scanCmd) scan(port int, c chan<- string) {
	if cmd.tCPonly {
		if cmd.shouldLog(tcp, port) {
			c <- fmt.Sprintf("%s\nport : %s\nOpen : %t",
				stars,
				fmt.Sprintf("%s/%d", tcp, port),
				isOpen(tcp, cmd.addr, port),
			)
		}
	}

	if cmd.uDPonly {
		if cmd.shouldLog(udp, port) {
			c <- fmt.Sprintf("%s\nport : %s\nOpen : %t",
				stars,
				fmt.Sprintf("%s/%d", udp, port),
				isOpen(udp, cmd.addr, port),
			)
		}
	}
}

func (cmd *scanCmd) shouldLog(protocol string, port int) bool {
	return cmd.openOnly && isOpen(protocol, cmd.addr, port) || !cmd.openOnly
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

func organize(results []string) []string {
	var organized []string
	for i := 0; i < len(results); i++ {
		for _, r := range results {
			if strings.Contains(r, strconv.Itoa(i)) {
				organized = append(organized, r)
			}
		}
	}
	return organized
}

func (cmd *scanCmd) start(portsForScanning []int) {
	var wg sync.WaitGroup
	var results []string
	wg.Add(len(portsForScanning))
	ch := make(chan string)

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

	flog.Info("starting scan...")

	for _, port := range portsForScanning {
		go func(p int) {
			cmd.scan(p, ch)
			wg.Done()
		}(port)
	}
	wg.Wait()
	close(ch)

	for _, r := range organize(results) {
		flog.Info(r)
	}
}
