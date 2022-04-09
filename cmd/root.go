package cmd

import (
	_ "embed"

	"github.com/spf13/cobra"
)

var (
	//go:embed version.txt
	version             string
	output              string
	shouldOutputVersion bool
)

func init() {
	Root.PersistentFlags().StringVarP(&output, "output", "o", output, "Output format. Supported values include json and yaml.")
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
