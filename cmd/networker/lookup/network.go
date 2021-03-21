package lookup

import (
	"log"

	"github.com/spf13/pflag"
	"go.coder.com/cli"

	"github.com/fuskovic/networker/internal/resolve"
)

type networkCmd struct {
	host string
}

func (cmd *networkCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "network",
		Usage: "[flags]",
		Desc:  "Lookup the network address of a provided host.",
	}
}

func (cmd *networkCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVar(&cmd.host, "host", "", "IP address or hostname to get the network address for.")
}

func (cmd *networkCmd) Run(fl *pflag.FlagSet) {
	if cmd.host == "" {
		fl.Usage()
		log.Fatal("no host provided")
	}

	network, err := resolve.NetworkByHost(cmd.host)
	if err != nil {
		fl.Usage()
		log.Fatalf("lookup failed: %s", err)
	}
	log.Printf("lookup successful - network address: %s", network)
}
