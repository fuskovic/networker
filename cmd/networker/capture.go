package main

import (
	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"

	cap "github.com/fuskovic/networker/internal/capture"
)

type captureCmd struct {
	out  string
	wide bool
}

// Spec returns a command spec containing a description of it's usage.
func (c *captureCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "capture",
		Usage:   "[flags]",
		Aliases: []string{"c", "cap"},
		Desc:    "Monitor network traffic on the LAN.",
	}
}

// RegisterFlags initializes how a flag set is processed for a particular command.
func (c *captureCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVarP(&c.out, "out", "o", c.out, "Name of an output file to write the packets to.")
	fl.BoolVarP(&c.wide, "wide", "w", c.wide, "Include hostnames, sequence, and mac addresses in output.")
}

// Run validates the flagset and runs the packet capture session accordingly.
func (c *captureCmd) Run(fl *pflag.FlagSet) {
	if err := c.capture(); err != nil {
		flog.Error("error running capture : %v", err)
		fl.Usage()
	}
}

func (c *captureCmd) capture() error {
	s := cap.Sniffer{
		File: c.out,
		Wide: c.wide,
	}
	return s.Capture()
}
