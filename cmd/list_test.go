package cmd

import (
	"encoding/json"
	"os/exec"

	"testing"

	"github.com/fuskovic/networker/v2/internal/list"
	"github.com/fuskovic/networker/v2/internal/test"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestListCommand(t *testing.T) {
	test.WithNetworker(t, "list devices output as json", func(t *testing.T) {
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
	test.WithNetworker(t, "list devices output as yaml", func(t *testing.T) {
		// start the list command
		cmd := exec.Command("networker", "ls", "-o", "yaml")
		stdout, err := cmd.StdoutPipe()
		require.NoError(t, err)
		require.NoError(t, cmd.Start())

		// assert we can unmarshal the yaml output as expected
		var devices []list.Device
		require.NoError(t, yaml.NewDecoder(stdout).Decode(&devices))
		require.NoError(t, cmd.Wait())

		// assert that the devices are not empty
		require.True(t, len(devices) > 0)
	})
}
