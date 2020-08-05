package main

import (
	"context"

	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"

	"github.com/fuskovic/networker/internal/list"
)

type listCmd struct{}

// Spec returns a command spec containing a description of it's usage.
func (c *listCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "list",
		Usage:   "[flags]",
		Aliases: []string{"ls"},
		Desc:    "List information on connected network devices.",
	}
}

// Run lists network devices.
func (c *listCmd) Run(fl *pflag.FlagSet) {
	err := list.List(context.Background())
	if err != nil {
		flog.Error("failed to list : %v", err)
		fl.Usage()
	}
}
