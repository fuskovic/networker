package networker

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
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "no ip address provided")
		})
		test.WithNetworker(t, "lookup hostname invalid ip address", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "hostname", "--ip", "invalid")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "not a valid ip address")
		})
		test.WithNetworker(t, "lookup ip address hostname not provided", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "ip")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "hostname not provided")
		})
		test.WithNetworker(t, "lookup isp no host provided", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "isp")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "no host provided")
		})
		test.WithNetworker(t, "lookup isp for private ip address", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "isp", "--host", "127.0.0.1")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "cannot retrieve internet service provider for private ip")
		})
		test.WithNetworker(t, "lookup nameservers hostname not provided", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "nameservers")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "hostname not provided")
		})
		test.WithNetworker(t, "lookup nameservers hostname doesnt exist", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "nameservers", "--hostname", "doesntexist")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "lookup failed")
		})
		test.WithNetworker(t, "lookup network no host provided", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "network")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "no host provided")
		})
	})
	t.Run("ShouldPass", func(t *testing.T) {
		test.WithNetworker(t, "lookup hostname", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "hostname", "--ip", "8.8.8.8")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "lookup successful")
		})
		test.WithNetworker(t, "lookup ip", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "ip", "--hostname", "dns.google.")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "lookup successful")
		})
		test.WithNetworker(t, "lookup isp with hostname", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "isp", "--host", "dns.google.")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "GOOGLE, US")
		})
		test.WithNetworker(t, "lookup isp with hostname as json output", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "isp", "--host", "dns.google.", "--json")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			isp := new(resolve.InternetServiceProvider)
			require.NoError(t, json.Unmarshal(output, isp))
		})
		test.WithNetworker(t, "lookup isp with ip address", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "isp", "--host", "8.8.8.8")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "GOOGLE, US")
		})
		test.WithNetworker(t, "lookup isp with ip address as json output", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "isp", "--host", "8.8.8.8", "--json")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			isp := new(resolve.InternetServiceProvider)
			require.NoError(t, json.Unmarshal(output, isp))
		})
	})
}
