package networker

import (
	"github.com/spf13/pflag"
	"go.coder.com/cli"

	"github.com/fuskovic/networker/cmd/networker/lookup"
)

type lookupCmd struct{}

func (cmd *lookupCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "lookup",
		Usage:   "[flags]",
		Aliases: []string{"lu"},
		Desc:    "Lookup hostnames, IP addresses, internet service providers, nameservers, and networks.",
	}
}

func (cmd *lookupCmd) Subcommands() []cli.Command {
	return []cli.Command{
		new(lookup.HostnameCmd),
		new(lookup.IpaddressCmd),
		new(lookup.NetworkCmd),
		new(lookup.NameserversCmd),
		new(lookup.IspCmd),
	}
}

func (cmd *lookupCmd) Run(fl *pflag.FlagSet) { fl.Usage() }
