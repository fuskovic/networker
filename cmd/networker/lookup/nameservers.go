package lookup

import (
	"encoding/json"
	"log"
	"os"

	"cdr.dev/coder-cli/pkg/tablewriter"
	"github.com/fuskovic/networker/internal/resolve"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
)

type nameserversCmd struct {
	hostname string
	json     bool
}

func (cmd *nameserversCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "nameservers",
		Usage: "[flags]",
		Desc:  "Lookup nameservers for the provided hostname.",
	}
}

func (cmd *nameserversCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVar(&cmd.hostname, "host", "", "Hostname to lookup nameservers for.")
}

func (cmd *nameserversCmd) Run(fl *pflag.FlagSet) {
	if cmd.hostname == "" {
		fl.Usage()
		log.Fatal("hostname not provided")
	}

	nameservers, err := resolve.NameServersByHostName(cmd.hostname)
	if err != nil {
		fl.Usage()
		log.Fatalf("lookup failed: %s", err)
	}

	if cmd.json {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "\t")
		enc.SetEscapeHTML(false)
		if err := enc.Encode(nameservers); err != nil {
			fl.Usage()
			log.Fatalf("failed to encode nameservers as json: %s", err)
		}
		return
	}

	if err := tablewriter.WriteTable(os.Stdout, len(nameservers), func(i int) interface{} { return nameservers[i] }); err != nil {
		fl.Usage()
		log.Fatalf("failed to write nameservers table: %s", err)
	}
}
