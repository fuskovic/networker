package main

import (
	"github.com/spf13/pflag"
	"go.coder.com/cli"
)

type lookupCmd struct {
	hostname   string
	ipAddress  string
	nameserver string
	network    string
}

func (cmd *lookupCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "lookup",
		Usage:   "[flags]",
		Aliases: []string{"lu"},
		Desc:    "Lookup hostnames, IP addresses, nameservers, and networks.",
	}
}

func (cmd *lookupCmd) Subcommands() []cli.Command {
	return []cli.Command{
		new(hostnameCmd),
		new(ipaddressCmd),
		new(networkCmd),
		new(nameserversCmd),
	}
}

func (cmd *lookupCmd) Run(fl *pflag.FlagSet) { fl.Usage() }
