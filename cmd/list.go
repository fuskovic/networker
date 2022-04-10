package cmd

import (
	"context"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/fuskovic/networker/v2/internal/encoder"
	"github.com/fuskovic/networker/v2/internal/list"
	"github.com/fuskovic/networker/v2/internal/usage"
)

func init() {
	Root.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List information on connected network devices.",
	Example: `
# List devices on network:

	networker list

# List devices on network(short-hand):

	nw ls

# List devices on network(short-hand) and output as json:

	nw ls -o json

# List devices on network(short-hand) and output as yaml:

	nw ls -o yaml
`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		s := spinner.New(spinner.CharSets[36], 50*time.Millisecond)
		s.Start()

		devices, err := list.Devices(ctx)
		if err != nil {
			usage.Fatalf(cmd, "failed to list devices: %s", err)
		}

		s.Stop()

		enc := encoder.New[list.Device](os.Stdout, output)
		if err := enc.Encode(devices...); err != nil {
			usage.Fatalf(cmd, "failed to encode devices: %s", err)
		}
	},
}
