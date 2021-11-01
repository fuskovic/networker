package networker

import (
	"github.com/spf13/pflag"
	"go.coder.com/cli"
)

type versionCmd struct{}

func (cmd *versionCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "version",
		Usage:   "[flags]",
		Aliases: []string{"v"},
		Desc:    "Print networker version.",
	}
}

func (cmd *versionCmd) Run(_ *pflag.FlagSet) { println("v1.2.4") }
