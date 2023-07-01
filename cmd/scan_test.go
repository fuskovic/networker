package cmd

import (
	"encoding/json"
	"os/exec"
	"testing"

	"github.com/fuskovic/networker/v3/internal/scanner"
	"github.com/fuskovic/networker/v3/internal/test"
	"github.com/stretchr/testify/require"
)

func TestScanCommand(t *testing.T) {
	test.WithNetworker(t, "output scanned devices as json", func(t *testing.T) {
		// start the list command
		cmd := exec.Command("networker", "scan", "-o", "json")
		stdout, err := cmd.StdoutPipe()
		require.NoError(t, err)
		require.NoError(t, cmd.Start())

		// assert we can unmarshal the json output as expected
		var scanResults []scanner.Scan
		require.NoError(t, json.NewDecoder(stdout).Decode(&scanResults))
		require.NoError(t, cmd.Wait())

		// assert that the results are not empty
		require.True(t, len(scanResults) > 0)
	})
	test.WithNetworker(t, "output scanned devices as yaml", func(t *testing.T) {
		// start the list command
		cmd := exec.Command("networker", "scan", "-o", "json")
		stdout, err := cmd.StdoutPipe()
		require.NoError(t, err)
		require.NoError(t, cmd.Start())

		// assert we can unmarshal the json output as expected
		var scanResults []scanner.Scan
		require.NoError(t, json.NewDecoder(stdout).Decode(&scanResults))
		require.NoError(t, cmd.Wait())

		// assert that the results are not empty
		require.True(t, len(scanResults) > 0)
	})
}
