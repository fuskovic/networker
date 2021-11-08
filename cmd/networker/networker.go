package networker

import (
	"github.com/spf13/pflag"
	"go.coder.com/cli"
)

type rootCmd struct{
	version bool
}

func (cmd *rootCmd) Run(fl *pflag.FlagSet) {
	if cmd.version {
		println("v1.2.5") 
		return
	}
	fl.Usage()
}

func (cmd *rootCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.BoolVarP(&cmd.version, "versin", "v", cmd.version, "Installed version.")
}

func (cmd *rootCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "networker",
		Usage: "[subcommand] [flags]",
		Desc:  "A simple networking tool.",
	}
}

func (cmd *rootCmd) Subcommands() []cli.Command {
	return []cli.Command{
		new(listCmd),
		new(lookupCmd),
		new(requestCmd),
		new(scanCmd),
	}
}

// Execute the root command.
func Execute() {
	cli.RunRoot(new(rootCmd))
}
