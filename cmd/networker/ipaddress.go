package main

import (
	"context"
	"os"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"

	"github.com/fuskovic/networker/internal/lookup"
)

type ipaddressCmd struct {
	hostname string
}

func (cmd *ipaddressCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "ip",
		Usage: "[flags]",
		Desc:  "Lookup the ip address of the provided hostname.",
	}
}

func (cmd *ipaddressCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVar(&cmd.hostname, "hostname", "", "Hostname to get the ip address of.")
}

func (cmd *ipaddressCmd) Run(fl *pflag.FlagSet) {
	if cmd.hostname == "" {
		fl.Usage()
		flog.Error("hostname not provided")
		return
	}
	ipAddr, err := lookup.AddrByHostName(cmd.hostname)
	if err != nil {
		fl.Usage()
		flog.Error("lookup failed: %v", err)
		return
	}
	sloghuman.Make(os.Stdout).Info(context.Background(), "lookup successful", slog.F("ip-address", ipAddr.String()))
}
