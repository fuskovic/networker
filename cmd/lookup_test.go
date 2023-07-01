package cmd

import (
	"encoding/json"
	"os/exec"
	"testing"

	"github.com/fuskovic/networker/v3/internal/resolve"
	"github.com/fuskovic/networker/v3/internal/test"
	"github.com/stretchr/testify/require"
)

func TestLookupHostnameCommand(t *testing.T) {
	t.Run("ShouldPass", func(t *testing.T) {
		test.WithNetworker(t, "lookup hostname", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "hostname", "8.8.8.8")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "dns.google.")
		})
		test.WithNetworker(t, "lookup hostname output as json", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "hostname", "8.8.8.8", "-o", "json")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			record := new(resolve.Record)
			require.NoError(t, json.Unmarshal(output, record))
			require.Equal(t, "dns.google.", record.Hostname)

		})
		test.WithNetworker(t, "lookup hostname output as yaml", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "hostname", "8.8.8.8", "-o", "yaml")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "- hostname: dns.google.")
		})
	})
	t.Run("ShouldFail", func(t *testing.T) {
		test.WithNetworker(t, "lookup hostname no ip address provided", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "hostname")
			output, _ := cmd.CombinedOutput()
			require.Contains(t, string(output), "Error: accepts 1 arg(s), received 0")
		})
		test.WithNetworker(t, "lookup hostname invalid ip address", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "hostname", "invalid")
			output, _ := cmd.CombinedOutput()
			require.Contains(t, string(output), "not a valid ip address")
		})
	})
}

func TestLookupIpAddressCommand(t *testing.T) {
	t.Run("ShouldPass", func(t *testing.T) {
		test.WithNetworker(t, "lookup ip", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "ip", "dns.google.")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "8.8.")
		})
		test.WithNetworker(t, "lookup ip output as json", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "ip", "dns.google.", "-o", "json")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			record := new(resolve.Record)
			require.NoError(t, json.Unmarshal(output, record))
			require.Equal(t, "dns.google.", record.Hostname)

		})
		test.WithNetworker(t, "lookup ip output as yaml", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "ip", "dns.google.", "-o", "yaml")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "- hostname: dns.google.")
		})
	})
	t.Run("ShouldFail", func(t *testing.T) {
		test.WithNetworker(t, "lookup ip address hostname not provided", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "ip")
			output, _ := cmd.CombinedOutput()
			require.Contains(t, string(output), "Error: accepts 1 arg(s), received 0")
		})
		test.WithNetworker(t, "lookup ip address but ip provided instead of hostname", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "ip", "8.8.8.8")
			output, _ := cmd.CombinedOutput()
			require.Contains(t, string(output), `expected a hostname not an ip address`)
		})
	})
}

func TestLookupIspCommand(t *testing.T) {
	t.Run("ShouldPass", func(t *testing.T) {
		test.WithNetworker(t, "lookup isp with hostname", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "isp", "dns.google.")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "GOOGLE, US")
		})
		test.WithNetworker(t, "lookup isp with hostname output as json output", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "isp", "dns.google.", "-o", "json")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			isp := new(resolve.InternetServiceProvider)
			require.NoError(t, json.Unmarshal(output, isp))
		})
		test.WithNetworker(t, "lookup isp with hostname output as yaml output", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "isp", "dns.google.", "-o", "yaml")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "- name: GOOGLE, US")
		})
		test.WithNetworker(t, "lookup isp with ip address", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "isp", "8.8.8.8")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "GOOGLE, US")
		})
		test.WithNetworker(t, "lookup isp with ip address output as json output", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "isp", "8.8.8.8", "-o", "json")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			isp := new(resolve.InternetServiceProvider)
			require.NoError(t, json.Unmarshal(output, isp))
		})
		test.WithNetworker(t, "lookup isp with ip address output as yaml output", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "isp", "8.8.8.8", "-o", "yaml")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "- name: GOOGLE, US")
		})
	})
	t.Run("ShouldFail", func(t *testing.T) {
		test.WithNetworker(t, "lookup isp no host provided", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "isp")
			output, _ := cmd.CombinedOutput()
			require.Contains(t, string(output), "Error: accepts 1 arg(s), received 0")
		})
		test.WithNetworker(t, "lookup isp for private ip address", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "isp", "127.0.0.1")
			output, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(output), "cannot retrieve internet service provider for private ip")
		})
	})
}

func TestLookupNetworkCommand(t *testing.T) {
	t.Run("ShouldPass", func(t *testing.T) {
		test.WithNetworker(t, "lookup network with hostname", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "network", "dns.google.")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "8.0.0.0")
		})
		test.WithNetworker(t, "lookup network with hostname output as json", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "network", "dns.google.", "-o", "json")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			record := new(resolve.NetworkRecord)
			require.NoError(t, err, json.Unmarshal(output, &record))
			require.Equal(t, "8.0.0.0", record.NetworkIP.String())
		})
		test.WithNetworker(t, "lookup network with hostname output as yaml", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "network", "dns.google.", "-o", "yaml")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "network: 8.0.0.0")
		})
		test.WithNetworker(t, "lookup network with ip", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "network", "8.8.8.8")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "8.0.0.0")
		})
		test.WithNetworker(t, "lookup network with ip output as json", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "network", "8.8.8.8", "-o", "json")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			record := new(resolve.NetworkRecord)
			require.NoError(t, err, json.Unmarshal(output, &record))
			require.Equal(t, "8.0.0.0", record.NetworkIP.String())
		})
		test.WithNetworker(t, "lookup network with ip output as yaml", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "network", "8.8.8.8", "-o", "yaml")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "network: 8.0.0.0")
		})
	})
	t.Run("ShouldFail", func(t *testing.T) {
		test.WithNetworker(t, "lookup network no arg provided", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "network")
			output, _ := cmd.CombinedOutput()
			require.Contains(t, string(output), "Error: accepts 1 arg(s), received 0")
		})
		test.WithNetworker(t, "lookup network invalid arg", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "network", "invalid")
			output, _ := cmd.CombinedOutput()
			require.Contains(t, string(output), "invalild host")
		})
	})
}

func TestLookupNameserversCommand(t *testing.T) {
	t.Run("ShouldPass", func(t *testing.T) {
		test.WithNetworker(t, "lookup nameservers with hostname", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "nameservers", "dns.google.")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "ns1.zdns.google.")
		})
		test.WithNetworker(t, "lookup nameservers with hostname output as json", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "nameservers", "dns.google.", "-o", "json")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			var nameservers []resolve.NameServer
			require.NoError(t, err, json.Unmarshal(output, &nameservers))
			require.False(t, len(nameservers) == 0)
		})
		test.WithNetworker(t, "lookup nameservers with hostname output as yaml", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "nameservers", "dns.google.", "-o", "yaml")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "host: ns1.zdns.google.")
		})
		test.WithNetworker(t, "lookup nameservers with ip", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "nameservers", "8.8.8.8")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "ns1.zdns.google.")
		})
		test.WithNetworker(t, "lookup nameservers with ip output as json", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "nameservers", "8.8.8.8", "-o", "json")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			var nameservers []resolve.NameServer
			require.NoError(t, err, json.Unmarshal(output, &nameservers))
			require.False(t, len(nameservers) == 0)
		})
		test.WithNetworker(t, "lookup nameservers with ip output as yaml", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "nameservers", "8.8.8.8", "-o", "yaml")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Contains(t, string(output), "host: ns1.zdns.google.")
		})
	})
	t.Run("ShouldFail", func(t *testing.T) {
		test.WithNetworker(t, "lookup nameservers hostname not provided", func(t *testing.T) {
			cmd := exec.Command("networker", "lookup", "nameservers")
			output, _ := cmd.CombinedOutput()
			require.Contains(t, string(output), "Error: accepts 1 arg(s), received 0")
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
			require.Contains(t, string(output), "Error: accepts 1 arg(s), received 0")
		})
	})
}
