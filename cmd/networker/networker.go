package networker

import (
	"github.com/spf13/pflag"
	"go.coder.com/cli"
)

type rootCmd struct{}

func (r *rootCmd) Run(fl *pflag.FlagSet) { fl.Usage() }

func (r *rootCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "networker",
		Usage: "[subcommand] [flags]",
		Desc:  "A simple networking tool.",
	}
}

func (r *rootCmd) Subcommands() []cli.Command {
	return []cli.Command{
		new(listCmd),
		new(lookupCmd),
		new(requestCmd),
		new(scanCmd),
		new(versionCmd),
	}
}

// Execute the root command.
func Execute() {
	cli.RunRoot(new(rootCmd))
}
