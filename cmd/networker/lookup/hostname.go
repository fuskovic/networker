package lookup

import (
	"context"
	"net"
	"os"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"github.com/fuskovic/networker/internal/resolve"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"
)

type hostnameCmd struct {
	ipAddress string
}

func (cmd *hostnameCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "hostname",
		Usage: "[flags]",
		Desc:  "Lookup the hostname for a provided ip address.",
	}
}

func (cmd *hostnameCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVar(&cmd.ipAddress, "ip", "", "IP address to get the hostname of.")
}

func (cmd *hostnameCmd) Run(fl *pflag.FlagSet) {
	if cmd.ipAddress == "" {
		fl.Usage()
		flog.Error("no ip address provided")
		return
	}

	ipAddr := net.ParseIP(cmd.ipAddress)
	if ipAddr == nil {
		fl.Usage()
		flog.Error("%q is not a valid ip address", cmd.ipAddress)
		return
	}

	hostname, err := resolve.HostNameByIP(ipAddr)
	if err != nil {
		fl.Usage()
		flog.Error("lookup failed: %v", err)
		return
	}
	slog.Make(sloghuman.Sink(os.Stdout)).Info(context.Background(), "lookup successful", slog.F("hostname", hostname))
}
