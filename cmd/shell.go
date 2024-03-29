package cmd

import (
	"net"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/fuskovic/networker/v3/internal/shell"
	"github.com/fuskovic/networker/v3/internal/usage"
)

var port int

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 4444, "Port to serve shell on.")
	shellCmd.AddCommand(serveCmd)
	shellCmd.AddCommand(dialCmd)
	Root.AddCommand(shellCmd)
}

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Serve and establish connections with remote shells.",
	Example: `
# Serve using the defaults(bash is the default shell and 4444 is the default port):

	nw shell serve

# Serve a particular shell on a particular port:

	nw shell serve zsh --port 9000

# Establish a new shell session by dialing a networker initiated shell server.

	nw shell dial some.remote.ip.addr:9000

`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Usage()
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a shell server.",
	Example: `
# Serve using the defaults(bash is the default shell and 4444 is the default port):

	nw shell serve

# Serve a particular shell on a particular port:

	nw shell serve zsh -p 9000

`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS == "windows" {
			usage.Fatal(cmd, "this command is not supported on windows")
		}

		sh := "bash"
		if len(args) == 1 {
			sh = strings.Replace(args[0], "/bin/", "", 1)
		}

		if err := shell.Serve(sh, port); err != nil {
			usage.Fatalf(cmd, "unexpected server shutdown: %s\n", err)
		}
	},
}

var dialCmd = &cobra.Command{
	Use:   "dial",
	Short: "Dial a shell server.",
	Example: `
# Establish a new shell session by dialing a networker initiated shell server.

	nw shell dial some.remote.ip.addr:9000

`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if _, _, err := net.SplitHostPort(args[0]); err != nil {
			usage.Fatalf(cmd, "invalid address: %s\n", err)
		}

		if err := shell.Dial(args[0]); err != nil {
			usage.Fatalf(cmd, "dialer error: %s\n", err)
		}
	},
}
