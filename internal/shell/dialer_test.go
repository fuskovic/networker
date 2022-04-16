package shell

import (
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDialer(t *testing.T) {
	t.Run("ShouldFail", func(t *testing.T) {
		t.Run("if addr is not serving shell", func(t *testing.T) {
			require.Error(t, Dial("localhost:80"))
		})
	})
	t.Run("ShouldPass", func(t *testing.T) {
		// get the process id of the current shell
		getShellPid := exec.Command("bash", "-c", "echo $$")
		output, err := getShellPid.CombinedOutput()
		require.NoError(t, err)
		out := strings.TrimSpace(string(output))
		ogPid, err := strconv.Atoi(out)
		require.NoError(t, err)

		go func() {
			// start the server
			require.NoError(t, Serve("bash", 4444))
		}()

		// Connect to it
		go func() {
			require.NoError(t, Dial("localhost:4444"))
		}()

		// grace period to wait for connection to establish
		time.Sleep(time.Second)

		// get the process id of the new shell
		getShellPid = exec.Command("bash", "-c", "echo $$")
		output, err = getShellPid.CombinedOutput()
		require.NoError(t, err)
		out = strings.TrimSpace(string(output))
		newPid, err := strconv.Atoi(out)
		require.NoError(t, err)

		// assert that the current shells process id is different than the original
		require.NotEqual(t, ogPid, newPid)
	})
}
