package cmd

import (
	"net"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/fuskovic/networker/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestServeCommand(t *testing.T) {
	t.Run("ShouldFail", func(t *testing.T) {
		test.WithNetworker(t, "if shell is unsupported", func(t *testing.T) {
			cmd := exec.Command("networker", "shell", "serve", "--shell", "unsupported")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "shell \"unsupported\" is not supported")
		})
		test.WithNetworker(t, "if shell is not installed", func(t *testing.T) {
			cmd := exec.Command("networker", "shell", "serve", "--shell", "fish")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t,
				string(output),
				"shell \"fish\" does not exist on system: exec: \"fish\": executable file not found in $PATH",
			)
		})
		test.WithNetworker(t, "if port is invalid", func(t *testing.T) {
			cmd := exec.Command("networker", "shell", "serve", "70000")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "\"70000\" is not a valid port")
		})
	})
	t.Run("ShouldPass", func(t *testing.T) {
		test.WithNetworker(t, "if args are valid", func(t *testing.T) {
			cmd := exec.Command("networker", "shell", "serve", "8000")
			require.NoError(t, cmd.Start())

			// validate that the server is up
			conn, err := net.DialTimeout("tcp", "localhost:8000", 3*time.Second)
			require.NoError(t, err)

			// close the client connection
			require.NoError(t, conn.Close())

			// kill the server
			require.NoError(t,
				syscall.Kill(syscall.Getpid(), syscall.SIGINT),
			)
		})
	})
}
