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

func (cmd *captureCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "capture",
		Usage:   "[flags]",
		Aliases: []string{"c", "cap"},
		Desc:    "Monitor network traffic on the LAN.",
	}
}

func (cmd *captureCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVarP(&cmd.out, "out", "o", cmd.out, "Name of an output file to write the packets to.")
	fl.BoolVarP(&cmd.wide, "wide", "w", cmd.wide, "Include hostnames, sequence number, and mac addresses in output.")
}

func (cmd *captureCmd) Run(fl *pflag.FlagSet) {
	if err := cmd.capture(); err != nil {
		fl.Usage()
		flog.Error("error running capture : %v", err)
	}
}

func (cmd *captureCmd) capture() error {
	return (&cap.Sniffer{
		File: cmd.out,
		Wide: cmd.wide,
	}).Run()
}
