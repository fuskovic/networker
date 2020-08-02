package main

import (
	"context"
	"fmt"
	"net"

	"github.com/fuskovic/networker/internal/scan"
	u "github.com/fuskovic/networker/internal/utils"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"
)

type scanCmd struct {
	host    string
	verbose bool
}

// Spec returns a command spec containing a description of it's usage.
func (c *scanCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "scan",
		Usage: "[flags]",
		Desc:  "Scan the well-known ports of a given host.",
	}
}

// RegisterFlags initializes how a flag set is processed for a particular command.
func (c *scanCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVar(&c.host, "host", "", "Host to scan.")
	fl.BoolVarP(&c.verbose, "verbose", "v", c.verbose, "Stream live scan results.")
}

// Run figures out which ports to scan using the flagset and scans them.
func (c *scanCmd) Run(fl *pflag.FlagSet) {
	ctx := context.Background()
	if err := c.validate(); err != nil {
		flog.Error("failed validation : %v", err)
		fl.Usage()
		return
	}
	c.scan(ctx)
}

func (c *scanCmd) validate() error {
	if net.ParseIP(c.host) == nil {
		if u.AddrByHostName(c.host) == "" {
			return fmt.Errorf("'%s' is not a valid host.", c.host)
		}
	}
	return nil
}

func (c *scanCmd) scan(ctx context.Context) {
	s := &scan.Scanner{
		Host:    c.host,
		Ports:   portsToScan(1024),
		Verbose: c.verbose,
	}
	s.Scan(ctx)
}

func portsToScan(max int) []int {
	var ports []int
	for p := 0; p < max; p++ {
		ports = append(ports, p)
	}
	return ports
}
