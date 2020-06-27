package cmd

import (
	"go.coder.com/cli"

	"github.com/spf13/pflag"
)

// RootCmd is the command that starts the program.
type RootCmd struct{}

// Run prints the usage of a flag set.
func (r *RootCmd) Run(fl *pflag.FlagSet) { fl.Usage() }

// Spec returns a command spec containing a description of it's usage.
func (r *RootCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "networker",
		Usage: "[subcommand] [flags]",
		Desc:  "A practical CLI tool for network administration.",
	}
}

// Subcommands returns a set of any existing child-commands.
func (r *RootCmd) Subcommands() []cli.Command {
	return []cli.Command{
		// TODO : ADD SUB-CMDS
	}
}
