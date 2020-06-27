package cmd

import (
	"go.coder.com/cli"

	"github.com/spf13/pflag"
)

// Root is the command that starts the program.
type Root struct{}

// Run prints the usage of a flag set.
func (r *Root) Run(fl *pflag.FlagSet) { fl.Usage() }

// Spec returns a command spec containing a description of it's usage.
func (r *Root) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "networker",
		Usage: "[subcommand] [flags]",
		Desc:  "A practical CLI tool for network administration.",
	}
}

// Subcommands returns a set of any existing child-commands.
func (r *Root) Subcommands() []cli.Command {
	return []cli.Command{
		&backDoorCmd{},
	}
}
