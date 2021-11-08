package lookup

import (
	"encoding/json"
	"os"

	"cdr.dev/coder-cli/pkg/tablewriter"
	"github.com/fuskovic/networker/internal/resolve"
	"github.com/fuskovic/networker/internal/usage"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
)

type NameserversCmd struct {
	hostname string
	json     bool
}

func (cmd *NameserversCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "nameservers",
		Usage: "[flags]",
		Desc:  "Lookup nameservers for the provided hostname.",
	}
}

func (cmd *NameserversCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVar(&cmd.hostname, "hostname", "", "Hostname to lookup nameservers for.")
}

func (cmd *NameserversCmd) Run(fl *pflag.FlagSet) {
	if cmd.hostname == "" {
		usage.Fatal(fl, "hostname not provided")
	}

	nameservers, err := resolve.NameServersByHostName(cmd.hostname)
	if err != nil {
		usage.Fatalf(fl, "lookup failed: %s", err)
	}

	if cmd.json {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "\t")
		enc.SetEscapeHTML(false)
		if err := enc.Encode(nameservers); err != nil {
			usage.Fatalf(fl, "failed to encode nameservers as json: %s", err)
		}
		return
	}

	err = tablewriter.WriteTable(os.Stdout, len(nameservers),
		func(i int) interface{} {
			return nameservers[i]
		},
	)

	if err != nil {
		usage.Fatalf(fl, "failed to write nameservers table: %s", err)
	}
}
