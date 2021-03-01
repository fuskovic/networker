package main

import (
	"context"
	"os"

	"cdr.dev/slog/sloggers/sloghuman"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"

	"github.com/fuskovic/networker/internal/list"
)

type listCmd struct{}

func (c *listCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "list",
		Usage:   "[flags]",
		Aliases: []string{"ls"},
		Desc:    "List information on connected network devices.",
	}
}

func (c *listCmd) Run(fl *pflag.FlagSet) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	devices, err := list.Devices(ctx)
	if err != nil {
		fl.Usage()
		flog.Error("failed to list devices: %v", err)
		return
	}

	log := sloghuman.Make(os.Stdout)
	for i := range devices {
		log.Info(ctx, string(devices[i].Kind()), devices[i].Fields()...)
	}
}
