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

	var hosts []string

	if cmd.host == "" {
		devices, err := list.Devices(ctx)
		if err != nil {
			fl.Usage()
			flog.Error("failed to list network devices : %v", err)
			return
		}
		for i := range devices {
			hosts = append(hosts, devices[i].Addr())
		}
	} else {
		ip := net.ParseIP(cmd.host)
		if ip == nil {
			if _, err := lookup.AddrByHostName(cmd.host); err != nil {
				fl.Usage()
				flog.Error("%q is an invalid host : %v", cmd.host, err)
				return
			}
		}
		hosts = append(hosts, cmd.host)
	}

	log := sloghuman.Make(os.Stdout)
	log.Info(ctx, "scanning", slog.F("hosts", hosts))

	for host, openPorts := range ports.NewScanner(hosts, cmd.all).Scan(ctx) {
		var hostname string

		ip := net.ParseIP(host)
		if ip == nil {
			hostname = host
			addr, err := lookup.AddrByHostName(hostname)
			if err != nil {
				flog.Error("failed to lookup ip address for hostname %q: %w", hostname, err)
				continue
			}
			ip = *addr
		} else {
			host, err := lookup.HostNameByIP(ip)
			if err != nil {
				flog.Error("failed to lookup hostname for ip %q: %w", ip, err)
				continue
			}
			hostname = host
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
