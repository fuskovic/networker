package lookup

import (
	"log"
	"net"

	"github.com/fuskovic/networker/internal/resolve"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
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
		log.Fatal("no ip address provided")
	}

	ipAddr := net.ParseIP(cmd.ipAddress)
	if ipAddr == nil {
		fl.Usage()
		log.Fatalf("%q is not a valid ip address", cmd.ipAddress)
	}

	hostname, err := resolve.HostNameByIP(ipAddr)
	if err != nil {
		fl.Usage()
		log.Fatalf("lookup failed: %s", err)
		return
	}
	log.Printf("lookup successful - hostname: %s", hostname)
}
