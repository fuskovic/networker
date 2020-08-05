package main

import (
	"github.com/spf13/pflag"
	"go.coder.com/cli"
)

type root struct{}

// Run prints the usage of a flag set.
func (r *root) Run(fl *pflag.FlagSet) { fl.Usage() }

// Spec returns a command spec containing a description of it's usage.
func (r *root) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "networker",
		Usage: "[subcommand] [flags]",
		Desc:  "A practical CLI tool for network administration.",
	}
}

// Subcommands returns a set of any existing child-commands.
func (r *root) Subcommands() []cli.Command {
	return []cli.Command{
		&captureCmd{},
		&listCmd{},
		&lookUpCmd{},
		&requestCmd{},
		&scanCmd{},
	}
}
