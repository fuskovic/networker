package main

import (
	"context"
	"net"
	"os"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"github.com/fuskovic/networker/internal/list"
	"github.com/fuskovic/networker/internal/lookup"
	"github.com/fuskovic/networker/internal/ports"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"
	"golang.org/x/xerrors"
)

type scanCmd struct {
	host string
	all  bool
}

func (cmd *scanCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "scan",
		Usage:   "[flags]",
		Aliases: []string{"s"},
		Desc:    "Scan the well-known ports of a given host or network.",
	}
}

func (cmd *scanCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVar(&cmd.host, "host", "", "Host to scan.")
	fl.BoolVarP(&cmd.all, "all", "a", false, "Scan all ports(scans first 1024 if not enabled).")
}

func (cmd *scanCmd) Run(fl *pflag.FlagSet) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hosts, err := cmd.getHostsToScan(ctx)
	if err != nil {
		fl.Usage()
		flog.Error("failed to get hosts to scan: %w", err)
		return
	}

	log := slog.Make(sloghuman.Sink(os.Stdout))
	log.Info(ctx, "scanning", slog.F("hosts", hosts))

	for host, openPorts := range ports.NewScanner(hosts, cmd.all).Scan(ctx) {
		hostname, ip, err := resolveHostAndAddr(host)
		if err != nil {
			fl.Usage()
			flog.Error("failed to resolve host and ip address for %q: %w", host, err)
			return
		}

		var foundOpenPorts bool
		if len(openPorts) > 0 {
			foundOpenPorts = true
		}

		log.Info(ctx, "scan complete",
			slog.F("ip-address", ip.String()),
			slog.F("found-open-ports", foundOpenPorts),
			slog.F("hostname", hostname),
			slog.F("open-ports", openPorts),
		)
	}
}

// if the user did not provide a host, get all the hosts on the network.
func (cmd *scanCmd) getHostsToScan(ctx context.Context) ([]string, error) {
	if cmd.host == "" {
		var hosts []string
		devices, err := list.Devices(ctx)
		if err != nil {
			return nil, xerrors.Errorf("failed to list devices: %w", err)
		}
		for i := range devices {
			hosts = append(hosts, devices[i].Addr())
		}
		return hosts, nil
	}

	if !isValidHost(cmd.host) {
		return nil, xerrors.Errorf("%q is not a valid host", cmd.host)
	}
	return []string{cmd.host}, nil
}

func resolveHostAndAddr(host string) (string, *net.IP, error) {
	var hostname string

	ip := net.ParseIP(host)
	if ip == nil {
		hostname = host
		addr, err := lookup.AddrByHostName(hostname)
		if err != nil {
			return "", nil, xerrors.Errorf("failed to get ip address by hostname %q: %w", hostname, err)
		}
		ip = *addr
	} else {
		host, err := lookup.HostNameByIP(ip)
		if err != nil {
			return "", nil, xerrors.Errorf("failed to get hostname by ip address %q: %w", ip, err)
		}
		hostname = host
	}
	return hostname, &ip, nil
}

func isValidHost(host string) bool {
	// otherwise, the user may have either provided a hostname or an ip address.
	ip := net.ParseIP(host)
	//if the ip address parse failed, we can assume the user provided a hostname.
	if ip == nil {
		// if we can't look up an ip address by the provided hostname, we can assume the hostname provided was invalid.
		if _, err := lookup.AddrByHostName(host); err != nil {
			return false
		}
	}
	return true
}
