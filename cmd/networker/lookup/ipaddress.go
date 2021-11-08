package lookup

import (
	"log"

	"github.com/fuskovic/networker/internal/resolve"
	"github.com/fuskovic/networker/internal/usage"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
)

type IpaddressCmd struct {
	hostname string
}

func (cmd *IpaddressCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "ip",
		Usage: "[flags]",
		Desc:  "Lookup the ip address of the provided hostname.",
	}
}

func (cmd *IpaddressCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVar(&cmd.hostname, "hostname", "", "Hostname to get the ip address of.")
}

func (cmd *IpaddressCmd) Run(fl *pflag.FlagSet) {
	if cmd.hostname == "" {
		usage.Fatal(fl, "hostname not provided")
	}

	ipAddr, err := resolve.AddrByHostName(cmd.hostname)
	if err != nil {
		usage.Fatalf(fl, "lookup failed: %s", err)
	}
	log.Printf("lookup successful - ip-address: %s", ipAddr)
}
