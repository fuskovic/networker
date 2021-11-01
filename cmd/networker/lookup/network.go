package lookup

import (
	"log"

	"github.com/spf13/pflag"
	"go.coder.com/cli"

	"github.com/fuskovic/networker/internal/resolve"
	"github.com/fuskovic/networker/internal/usage"
)

type NetworkCmd struct {
	host string
}

func (cmd *NetworkCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "network",
		Usage: "[flags]",
		Desc:  "Lookup the network address of a provided host.",
	}
}

func (cmd *NetworkCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVar(&cmd.host, "host", "", "IP address or hostname to get the network address for.")
}

func (cmd *NetworkCmd) Run(fl *pflag.FlagSet) {
	if cmd.host == "" {
		usage.Fatal(fl, "no host provided")
	}

	network, err := resolve.NetworkByHost(cmd.host)
	if err != nil {
		usage.Fatalf(fl, "lookup failed: %s", err)
	}
	log.Printf("lookup successful - network address: %s", network)
}
