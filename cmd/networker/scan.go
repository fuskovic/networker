package main

import (
	"context"
	"net"
	"os"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"github.com/fuskovic/networker/internal/lookup"
	"github.com/fuskovic/networker/internal/ports"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"
)

type scanCmd struct {
	host string
}

func (cmd *scanCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "scan",
		Usage:   "[flags]",
		Aliases: []string{"s"},
		Desc:    "Scan the well-known ports of a given host.",
	}
}

func (cmd *scanCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVar(&cmd.host, "host", "", "Host to scan.")
}

func (cmd *scanCmd) Run(fl *pflag.FlagSet) {
	ip := net.ParseIP(cmd.host)
	if ip == nil {
		if _, err := lookup.AddrByHostName(cmd.host); err != nil {
			fl.Usage()
			flog.Error("%q is an invalid host : %v", cmd.host, err)
			return
		}
	}

	ctx := context.Background()
	openPorts := ports.NewScanner(cmd.host).Scan(ctx, portsToScan(1024))
	sloghuman.Make(os.Stdout).Info(ctx, "scan complete",
		slog.F("host", cmd.host),
		slog.F("open-ports", openPorts),
	)
}

func portsToScan(max int) []int {
	var ports []int
	for p := 0; p < max; p++ {
		ports = append(ports, p)
	}
	return ports
}
