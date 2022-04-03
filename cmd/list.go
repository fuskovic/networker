package cmd

import (
	"context"
	"encoding/json"
	"os"

	"cdr.dev/coder-cli/pkg/tablewriter"
	"github.com/spf13/cobra"

	"github.com/fuskovic/networker/internal/list"
	"github.com/fuskovic/networker/internal/usage"
)


func init() {
	ListCmd.Flags().BoolVar(&shouldOutputAsJSON, "json", false, "Output as JSON.")
	Root.AddCommand(ListCmd)
}

var ListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List information on connected network devices.",
	Example: `
List devices on network:
	networker ls

Output as JSON:
	networker ls --json
`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		devices, err := list.Devices(ctx)
		if err != nil {
			usage.Fatalf(cmd, "failed to list devices: %s", err)
		}

		if shouldOutputAsJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "\t")
			enc.SetEscapeHTML(false)
			if err := enc.Encode(devices); err != nil {
				usage.Fatalf(cmd, "failed to encode devices as json: %s", err)
			}
			return
		}

		err = tablewriter.WriteTable(os.Stdout, len(devices),
			func(i int) interface{} {
				return devices[i]
			},
		)

		if err != nil {
			usage.Fatalf(cmd, "failed to write devices table: %s", err)
		}
	},
}
