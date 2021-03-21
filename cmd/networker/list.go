package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"cdr.dev/coder-cli/pkg/tablewriter"
	"github.com/spf13/pflag"
	"go.coder.com/cli"

	"github.com/fuskovic/networker/internal/list"
)

type listCmd struct {
	json bool
}

func (c *listCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "list",
		Usage:   "[flags]",
		Aliases: []string{"ls"},
		Desc:    "List information on connected network devices.",
	}
}

func (cmd *listCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.BoolVar(&cmd.json, "json", false, "Output as json.")
}

func (cmd *listCmd) Run(fl *pflag.FlagSet) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	devices, err := list.Devices(ctx)
	if err != nil {
		fl.Usage()
		log.Fatalf("failed to list devices: %s", err)
	}

	if cmd.json {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "\t")
		enc.SetEscapeHTML(false)
		if err := enc.Encode(devices); err != nil {
			fl.Usage()
			log.Fatalf("failed to encode devices as json: %s", err)
		}
		return
	}

	if err := tablewriter.WriteTable(os.Stdout, len(devices), func(i int) interface{} { return devices[i] }); err != nil {
		fl.Usage()
		log.Fatalf("failed to write devices table: %s", err)
	}
}
