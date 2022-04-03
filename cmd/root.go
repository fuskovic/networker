package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

var shouldOutputVersion bool

// Overwritten using ldflags on build
var Version = "development"

func init() {
	Root.Flags().BoolVarP(&shouldOutputVersion, "version", "v", false, "Print installed version.")
}

var Root = &cobra.Command{
	Use:     "networker",
	Aliases: []string{"nw"},
	Short:   "A simple networking utility.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if shouldOutputVersion {
			println(strings.Split(Version, "-")[0])
			return
		}
		_ = cmd.Usage()
	},
}