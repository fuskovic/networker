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

func (cmd *lookupCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVarP(&cmd.network, "network", "n", "", "Look up the network address of the provided host.")
	fl.StringVar(&cmd.hostname, "ip", "", "Look up the IP address of the provided hostname.")
	fl.StringVarP(&cmd.nameserver, "nameservers", "s", "", "Look up nameservers of the provided hostname.")
	fl.StringVar(&cmd.ipAddress, "hostnames", "", "Look up the hostname for a provided ip address.")
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
