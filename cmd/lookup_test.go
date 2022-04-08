package cmd

import (
	"encoding/json"
	"os/exec"
	"testing"

	"github.com/fuskovic/networker/internal/resolve"
	"github.com/fuskovic/networker/internal/test"
	"github.com/stretchr/testify/require"
)

func TestLookupCommand(t *testing.T) {
	t.Run("ShouldFail", func(t *testing.T) {
		test.WithNetworker(t, "lookup hostname no ip address provided", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "hostname")
			output, _ := cmd.CombinedOutput()
			require.Contains(t, string(output), `Error: accepts 1 arg(s), received 0`)
		})
		test.WithNetworker(t, "lookup hostname invalid ip address", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "hostname", "invalid")
			output, _ := cmd.CombinedOutput()
			require.Contains(t, string(output), "not a valid ip address")
		})
		test.WithNetworker(t, "lookup ip address hostname not provided", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "ip")
			output, _ := cmd.CombinedOutput()
			require.Contains(t, string(output), `Error: accepts 1 arg(s), received 0`)
		})
		test.WithNetworker(t, "lookup isp no host provided", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "isp")
			output, _ := cmd.CombinedOutput()
			require.Contains(t, string(output), `Error: accepts 1 arg(s), received 0`)
		})
		test.WithNetworker(t, "lookup isp for private ip address", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "isp", "127.0.0.1")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "cannot retrieve internet service provider for private ip")
		})
		test.WithNetworker(t, "lookup nameservers hostname not provided", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "nameservers")
			output, _ := cmd.CombinedOutput()
			require.Contains(t, string(output), `Error: accepts 1 arg(s), received 0`)
		})
		test.WithNetworker(t, "lookup nameservers hostname doesnt exist", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "nameservers", "doesntexist")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "lookup failed")
		})
		test.WithNetworker(t, "lookup network no host provided", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "network")
			output, _ := cmd.CombinedOutput()
			require.Contains(t, string(output), `Error: accepts 1 arg(s), received 0`)
		})
	})
	t.Run("ShouldPass", func(t *testing.T) {
		test.WithNetworker(t, "lookup hostname", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "hostname", "8.8.8.8")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "dns.google.")
		})
		test.WithNetworker(t, "lookup ip", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "ip", "dns.google.")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "8.8.")
		})
		test.WithNetworker(t, "lookup isp with hostname", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "isp", "dns.google.")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "GOOGLE, US")
		})
		test.WithNetworker(t, "lookup isp with hostname as json output", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "isp", "dns.google.", "-o", "json")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			isp := new(resolve.InternetServiceProvider)
			require.NoError(t, json.Unmarshal(output, isp))
		})
		test.WithNetworker(t, "lookup isp with ip address", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "isp", "8.8.8.8")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "GOOGLE, US")
		})
		test.WithNetworker(t, "lookup isp with ip address as json output", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "isp", "8.8.8.8", "-o", "json")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			isp := new(resolve.InternetServiceProvider)
			require.NoError(t, json.Unmarshal(output, isp))
		})
	})
}
