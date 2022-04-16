package cmd

import (
	"github.com/fuskovic/networker/v2/internal/shell"
	"github.com/fuskovic/networker/v2/internal/usage"
	"github.com/spf13/cobra"
)

var targetShell string

func init() {
	serveCmd.Flags().StringVar(&targetShell, "shell", "bash", "Shell to serve. e.g. bash, zsh, sh, etc...")
	shellCmd.AddCommand(serveCmd)
	Root.AddCommand(shellCmd)
}

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Serve and establish connections with remote shells.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Usage()
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a shell server.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		port := "4444"
		if len(args) == 1 {
			port = args[0]
		}

		if err := shell.Serve(targetShell, port); err != nil {
			usage.Fatalf(cmd, "unexpected server shutdown: %s\n", err)
		}
	},
}
