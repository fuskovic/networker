package lookup

import (
	"encoding/json"
	"os"

	"cdr.dev/coder-cli/pkg/tablewriter"
	"github.com/spf13/pflag"
	"go.coder.com/cli"

	"github.com/fuskovic/networker/internal/resolve"
	"github.com/fuskovic/networker/internal/usage"
)

type IspCmd struct {
	host string
	json bool
}

func (cmd *IspCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "isp",
		Usage: "[flags]",
		Desc:  "Lookup the internet service provider of a remote host.",
	}
}

func (cmd *IspCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVar(&cmd.host, "host", "", "IP address or hostname to get the network address for.")
	fl.BoolVar(&cmd.json, "json", false, "Output as json.")
}

func (cmd *IspCmd) Run(fl *pflag.FlagSet) {
	if cmd.host == "" {
		usage.Fatal(fl, "no host provided")
	}

	_, ip, err := resolve.HostAndAddr(cmd.host)
	if err != nil {
		usage.Fatalf(fl, "%q is an invalid host: %s", cmd.host, err)
	}

	if resolve.IsPrivate(ip) {
		usage.Fatalf(fl, "%q is not a remote ip", ip)
	}

	isp, err := resolve.ServiceProvider(ip)
	if err != nil {

		usage.Fatalf(fl, "failed to resolve internet service provider for %q: %s", cmd.host, err)
	}

	if cmd.json {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "\t")
		enc.SetEscapeHTML(false)
		if err := enc.Encode(isp); err != nil {
			usage.Fatalf(fl, "failed to encode internet service provider as json: %s", err)
		}
		return
	}

	err = tablewriter.WriteTable(os.Stdout, 1,
		func(_ int) interface{} {
			return *isp
		},
	)

	if err != nil {
		usage.Fatalf(fl, "failed to write service provider table for %q: %s", cmd.host, err)
	}
}
