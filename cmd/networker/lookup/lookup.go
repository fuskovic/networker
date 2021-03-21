package lookup

import (
	"github.com/spf13/pflag"
	"go.coder.com/cli"
)

type Cmd struct{}

func (cmd *Cmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "lookup",
		Usage:   "[flags]",
		Aliases: []string{"lu"},
		Desc:    "Lookup hostnames, IP addresses, nameservers, and networks.",
	}
}

func (cmd *Cmd) Subcommands() []cli.Command {
	return []cli.Command{
		new(hostnameCmd),
		new(ipaddressCmd),
		new(networkCmd),
		new(nameserversCmd),
		new(ispCmd),
	}
}

func (cmd *Cmd) Run(fl *pflag.FlagSet) { fl.Usage() }
