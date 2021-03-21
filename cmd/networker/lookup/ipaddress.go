package lookup

import (
	"log"

	"github.com/fuskovic/networker/internal/resolve"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
)

type ipaddressCmd struct {
	hostname string
}

func (cmd *ipaddressCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "ip",
		Usage: "[flags]",
		Desc:  "Lookup the ip address of the provided hostname.",
	}
}

func (cmd *ipaddressCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVar(&cmd.hostname, "hostname", "", "Hostname to get the ip address of.")
}

func (cmd *ipaddressCmd) Run(fl *pflag.FlagSet) {
	if cmd.hostname == "" {
		fl.Usage()
		log.Fatal("hostname not provided")
	}

	ipAddr, err := resolve.AddrByHostName(cmd.hostname)
	if err != nil {
		fl.Usage()
		log.Fatalf("lookup failed: %s", err)
	}
	log.Printf("lookup successful - ip-address: %s", ipAddr)
}
