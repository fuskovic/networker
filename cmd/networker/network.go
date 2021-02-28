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

type networkCmd struct {
	host string
}

func (cmd *networkCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "network",
		Usage: "[flags]",
		Desc:  "Lookup the network address of a provided host.",
	}
}

func (cmd *networkCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVar(&cmd.host, "host", "", "IP address or hostname to get the network address for.")
}

func (cmd *networkCmd) Run(fl *pflag.FlagSet) {
	if cmd.host == "" {
		fl.Usage()
		flog.Error("no host provided")
		return
	}

	network, err := lookup.NetworkByHost(cmd.host)
	if err != nil {
		fl.Usage()
		flog.Error("lookup failed: %v", err)
		return
	}
	sloghuman.Make(os.Stdout).Info(context.Background(), "lookup successful", slog.F("network-address", network.String()))
}
