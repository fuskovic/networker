package cmd

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"

	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"
)

type proxyCmd struct {
	listenOn int
	upStream string
}

// Spec returns a command spec containing a description of it's usage.
func (cmd *proxyCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "proxy",
		Usage: "TODO: ADD USAGE",
		Desc:  "Forward network traffic from one network connection to another.",
	}
}

// RegisterFlags initializes how a flag set is processed for a particular command.
func (cmd *proxyCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVarP(&cmd.upStream, "upstream", "u", cmd.upStream, "Address of server to forward traffic to.")
	fl.IntVarP(&cmd.listenOn, "listen-on", "l", cmd.listenOn, "Port to listen on.")
}

// Run creates a TCP listener and forwards anything received on that connection to the dialed upstream connection.
func (cmd *proxyCmd) Run(fl *pflag.FlagSet) {
	port := fmt.Sprintf(":%d", cmd.listenOn)

	flog.Info("starting listener on %s...", port)

	lsnr, err := net.Listen(tcp, port)
	if err != nil {
		flog.Fatal(err.Error())
	}
	defer lsnr.Close()

	flog.Info("dialing %s", cmd.upStream)

	upStr, err := net.Dial(tcp, cmd.upStream)
	if err != nil {
		flog.Fatal(err.Error())
	}
	defer upStr.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)

	go func() {
		for {
			sig := <-c
			flog.Info("received termination signal : %v", sig)
			lsnr.Close()
		}
	}()

	flog.Info("proxy started")

	for {
		conn, err := lsnr.Accept()
		if err != nil {
			flog.Fatal(err.Error())
		}
		defer conn.Close()

		go io.Copy(conn, upStr)
		io.Copy(upStr, conn)
	}
}
