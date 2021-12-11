package networker

import (
	"log"

	"github.com/spf13/pflag"
	"go.coder.com/cli"

	"github.com/fuskovic/networker/internal/loadbalancer"
	"github.com/fuskovic/networker/internal/usage"
)

type balanceCmd struct {
	targets  []string
	strategy string
	tls      bool
	key      string
	port     string
}

func (cmd *balanceCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "balance",
		Usage:   "[flags]",
		Aliases: []string{"b"},
		Desc:    "Load balance HTTP/HTTPs traffic across multiple servers.",
	}
}

func (cmd *balanceCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringSliceVarP(&cmd.targets, "targets", "t", cmd.targets, "Addresses to proxy requests to(format: host:port,host:port,etc...).")
	fl.StringVarP(&cmd.strategy, "strategy", "s", cmd.strategy, "Load-balancing strategy.")
	fl.BoolVar(&cmd.tls, "tls", cmd.tls, "Enable TLS.")
	fl.StringVarP(&cmd.key, "key", "k", cmd.key, "Path to public key file or raw string literal.")
	fl.StringVarP(&cmd.port, "port", "p", cmd.port, "Port to run the reverse proxy on.")
}

func (cmd *balanceCmd) Run(fl *pflag.FlagSet) {
	lb, err := loadbalancer.New(
		&loadbalancer.Config{
			Targets:   cmd.targets,
			Strategy:  cmd.strategy,
			EnableTLS: cmd.tls,
			PublicKey: cmd.key,
			Port:      ":" + cmd.port,
		},
	)

	if err != nil {
		usage.Fatalf(fl, "failed to initialize new proxy: %s", err)
	}
	log.Printf("load balancer exited: %s\n", <-lb.Balance())
}
