package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"

	cap "github.com/fuskovic/networker/internal/capture"
)

type captureCmd struct {
	device       string
	seconds      int64
	out          string
	limit        bool
	numToCapture int64
	wide         bool
}

// Spec returns a command spec containing a description of it's usage.
func (c *captureCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "capture",
		Usage: "[flags]",
		Desc:  "Capture network packets on a given device.",
	}
}

// RegisterFlags initializes how a flag set is processed for a particular command.
func (c *captureCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.Int64VarP(&c.seconds, "seconds", "s", c.seconds, "Amount of seconds to run capture for.")
	fl.StringVarP(&c.device, "device", "d", c.device, "Device to capture packets on.")
	fl.StringVarP(&c.out, "out", "o", c.out, "Name of an output file to write the packets to.")
	fl.BoolVarP(&c.wide, "wide", "w", c.wide, "Include hostnames, sequence, and mac addresses in output.")
}

// Run validates the flagset and runs the packet capture session accordingly.
func (c *captureCmd) Run(fl *pflag.FlagSet) {
	var err error

	switch {
	case len(c.device) == 0:
		err = errors.New("no designated devices")
	case c.seconds < 5:
		err = fmt.Errorf("capture must be at least 5 seconds long - your input : %d", c.seconds)
	case c.limit && c.numToCapture < 1:
		err = errors.New("use of --limit flag without use of --num flag\nPlease specify number of packets to limit capture\nminimum is 1")
	default:
		err = c.capture()
	}

	if err != nil {
		flog.Error("error running capture : %v", err)
		fl.Usage()
	}
}

func (c *captureCmd) capture() error {
	s := cap.Sniffer{
		Device: c.device,
		Time:   time.Duration(c.seconds) * time.Second,
		File:   c.out,
		Wide:   c.wide,
	}
	return s.Capture()
}
