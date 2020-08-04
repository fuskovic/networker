package main

import (
	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"

	l "github.com/fuskovic/networker/internal/lookup"
)

type lookUpCmd struct {
	hostName   string
	ipAddress  string
	nameServer string
	network    string
}

// Spec returns a command spec containing a description of it's usage.
func (c *lookUpCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "lookup",
		Usage:   "[flags]",
		Aliases: []string{"lu"},
		Desc:    "Lookup hostnames, IP addresses, nameservers, and networks.",
	}
}

// RegisterFlags initializes how a flag set is processed for a particular command.
func (c *lookUpCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVarP(&c.network, "network", "n", "", "Look up the network a given hostname belongs to.")
	fl.StringVarP(&c.ipAddress, "addresses", "a", "", "Look up IP addresses for a given hostname.")
	fl.StringVarP(&c.nameServer, "nameservers", "s", "", "Look up nameservers for a given hostname.")
	fl.StringVar(&c.hostName, "hostnames", "", "Look up hostnames for a given IP address.")
}

// Run runs all enabled lookups.
func (c *lookUpCmd) Run(fl *pflag.FlagSet) {
	if !c.valid() {
		fl.Usage()
		return
	}

	for k, lu := range c.lookUps() {
		if k != "" {
			if err := lu(k); err != nil {
				flog.Error("errors running lookups : %v", err)
				fl.Usage()
			}
		}
	}
}

func (c *lookUpCmd) lookUps() map[string]l.LookUpFunc {
	return map[string]l.LookUpFunc{
		c.hostName:   l.HostNamesByIP,
		c.ipAddress:  l.AddrsByHostName,
		c.nameServer: l.NameServersByHostName,
		c.network:    l.NwByHostName,
	}
}

func (c *lookUpCmd) valid() bool {
	lookUps := []string{c.hostName,
		c.ipAddress,
		c.nameServer,
		c.network,
	}

	for _, lu := range lookUps {
		if lu != "" {
			return true
		}
	}
	return false
}
