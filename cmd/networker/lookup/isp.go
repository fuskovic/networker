package lookup

import (
	"encoding/json"
	"log"
	"os"

	"cdr.dev/coder-cli/pkg/tablewriter"
	"github.com/spf13/pflag"
	"go.coder.com/cli"

	"github.com/fuskovic/networker/internal/resolve"
)

type ispCmd struct {
	host string
	json bool
}

func (cmd *ispCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "isp",
		Usage: "[flags]",
		Desc:  "Lookup the internet service provider of a remote host.",
	}
}

func (cmd *ispCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVar(&cmd.host, "host", "", "IP address or hostname to get the network address for.")
	fl.BoolVar(&cmd.json, "json", false, "Output as json.")
}

func (cmd *ispCmd) Run(fl *pflag.FlagSet) {
	if cmd.host == "" {
		fl.Usage()
		log.Fatal("no host provided")
	}

	_, ip, err := resolve.HostAndAddr(cmd.host)
	if err != nil {
		fl.Usage()
		log.Fatalf("%q is an invalid host: %s", cmd.host, err)
	}

	if resolve.IsPrivate(ip) {
		fl.Usage()
		log.Fatalf("%q is not a remote ip", ip)
	}

	isp, err := resolve.ServiceProvider(ip)
	if err != nil {
		fl.Usage()
		log.Fatalf("failed to resolve internet service provider for %q: %s", cmd.host, err)
	}

	if cmd.json {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "\t")
		enc.SetEscapeHTML(false)
		if err := enc.Encode(isp); err != nil {
			fl.Usage()
			log.Fatalf("failed to encode internet service provider as json: %s", err)
		}
		return
	}

	if err := tablewriter.WriteTable(os.Stdout, 1, func(_ int) interface{} { return *isp }); err != nil {
		fl.Usage()
		log.Fatalf("failed to write service provider table for %q: %s", cmd.host, err)
	}
}
