package lookup

import (
	"log"
	"net"

	"github.com/fuskovic/networker/internal/resolve"
	"github.com/fuskovic/networker/internal/usage"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
)

type HostnameCmd struct {
	ipAddress string
}

func (cmd *HostnameCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "hostname",
		Usage: "[flags]",
		Desc:  "Lookup the hostname for a provided ip address.",
	}
}

func (cmd *HostnameCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVar(&cmd.ipAddress, "ip", "", "IP address to get the hostname of.")
}

func (cmd *HostnameCmd) Run(fl *pflag.FlagSet) {
	if cmd.ipAddress == "" {
		usage.Fatal(fl, "no ip address provided")
	}

	ipAddr := net.ParseIP(cmd.ipAddress)
	if ipAddr == nil {
		usage.Fatalf(fl, "%q is not a valid ip address", cmd.ipAddress)
	}

	hostname, err := resolve.HostNameByIP(ipAddr)
	if err != nil {
		usage.Fatalf(fl, "lookup failed: %s", err)
	}
	log.Printf("lookup successful - hostname: %s", hostname)
}
