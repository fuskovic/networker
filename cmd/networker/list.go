package main

import (
	"context"

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
	if err := list.Run(context.Background()); err != nil {
		fl.Usage()
		flog.Error("failed to list : %v", err)
	}
}
