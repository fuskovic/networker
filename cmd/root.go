package cmd

import (
	_ "embed"

	"github.com/spf13/cobra"
)

var shouldOutputVersion bool

//go:embed version.txt
var version string

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
			print(version)
			return
		}
		_ = cmd.Usage()
	},
}
