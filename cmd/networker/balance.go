package networker

import (
	"crypto/tls"
	"log"
	"net/http"

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
	cert     string
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
	fl.StringVarP(&cmd.key, "key", "k", cmd.key, "Path to public key file.")
	fl.StringVarP(&cmd.cert, "cert", "c", cmd.cert, "Path to TLS certificate.")
	fl.StringVarP(&cmd.port, "port", "p", cmd.port, "Port to run the reverse proxy on.")
}

func (cmd *balanceCmd) Run(fl *pflag.FlagSet) {
	if cmd.port == "" {
		usage.Fatal(fl, "port is unset")
	}

	cert, err := tls.LoadX509KeyPair(cmd.cert, cmd.key)
	if err != nil {
		usage.Fatalf(fl, "failed to load cert: %s", err)
	}

	port := ":" + cmd.port
	lb, err := loadbalancer.New(
		&loadbalancer.Config{
			Hosts:     cmd.targets,
			Strategy:  cmd.strategy,
			EnableTLS: cmd.tls,
			TlsCert:   cert,
		},
	)

	if err != nil {
		usage.Fatalf(fl, "failed to initialize new proxy: %s", err)
	}

	c := make(chan error, 1)
	if cmd.tls {
		c <- http.ListenAndServeTLS(port, cmd.cert, cmd.key, lb)
	} else {
		c <- http.ListenAndServe(port, lb)
	}
	log.Printf("load balancer shutting down: %s\n", <-c)
}
