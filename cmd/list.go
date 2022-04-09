package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/fuskovic/networker/internal/encoder"
	"github.com/fuskovic/networker/internal/list"
	"github.com/fuskovic/networker/internal/usage"
)

func init() {
	Root.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List information on connected network devices.",
	Example: `
	List devices on network:
		networker ls

	Output as JSON:
		networker ls -o json
`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		devices, err := list.Devices(ctx)
		if err != nil {
			usage.Fatalf(cmd, "failed to list devices: %s", err)
		}

		enc := encoder.New[list.Device](os.Stdout, output)
		if err := enc.Encode(devices...); err != nil {
			usage.Fatalf(cmd, "failed to encode devices: %s", err)
		}
	},
}
