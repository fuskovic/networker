package shell

import (
	"errors"
	"net"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	t.Run("ShouldFail", func(t *testing.T) {
		t.Run("if shell is unsupported", func(t *testing.T) {
			expected := errors.New("shell \"unsupported\" is not supported")
			got := Serve("unsupported", 80)
			require.Error(t, got)
			require.Equal(t, expected, got)
		})
		t.Run("if shell is not installed", func(t *testing.T) {
			expected := "shell \"fish\" does not exist on system: exec: \"fish\": executable file not found in $PATH"
			got := Serve("fish", 80)
			require.Error(t, got)
			require.Equal(t, expected, got.Error())
		})
		t.Run("if port is negative number", func(t *testing.T) {
			expected := "-1 is not a valid port"
			got := Serve("bash", -1)
			require.Error(t, got)
			require.Equal(t, expected, got.Error())
		})
		t.Run("if port is invalid", func(t *testing.T) {
			expected := "70000 is not a valid port"
			got := Serve("bash", 70000)
			require.Error(t, got)
			require.Equal(t, expected, got.Error())
		})
	})
	t.Run("ShouldPass", func(t *testing.T) {
		go func() {
			// start the server
			require.NoError(t, Serve("bash", 4444))
		}()

		// validate that its up
		conn, err := net.DialTimeout("tcp", "localhost:4444", 3*time.Second)
		require.NoError(t, err)

		// close the client connection
		require.NoError(t, conn.Close())

		// kill the server
		require.NoError(t,
			syscall.Kill(syscall.Getpid(), syscall.SIGINT),
		)
	})
}
