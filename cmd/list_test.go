package cmd

import (
	"encoding/json"
	"os/exec"

	"testing"

	"github.com/fuskovic/networker/internal/list"
	"github.com/fuskovic/networker/internal/test"
	"github.com/stretchr/testify/require"
)

func TestListCommand(t *testing.T) {
	test.WithNetworker(t, "list devices as json output", func(t *testing.T) {
		// start the list command
		cmd := exec.Command("networker", "ls", "-o", "json")
		stdout, err := cmd.StdoutPipe()
		require.NoError(t, err)
		require.NoError(t, cmd.Start())

		// assert we can unmarshal the json output as expected
		var devices []list.Device
		require.NoError(t, json.NewDecoder(stdout).Decode(&devices))
		require.NoError(t, cmd.Wait())

		// assert that the devices are not empty
		require.True(t, len(devices) > 0)
	})
}
