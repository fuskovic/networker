package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

var shouldOutputVersion bool

// Overwritten using ldflags on build
var Version = "development"

func init() {
	rootCmd.Flags().BoolVarP(&shouldOutputVersion, "version", "v", false, "Print installed version.")
}

var rootCmd = &cobra.Command{
	Use:     "networker",
	Aliases: []string{"nw"},
	Short:   "A simple networking utility.",
	Example: `
Print version:
	networker -v
`,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if shouldOutputVersion {
			println(strings.Split(Version, "-")[0])
			return
		}
		_ = cmd.Usage()
	},
}

func Execute() {
	_ = rootCmd.Execute()
}
