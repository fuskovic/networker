package resolve

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLookup(t *testing.T) {
	t.Parallel()
	googleDNS := net.ParseIP("8.8.8.8")
	t.Run("ShouldPass", func(t *testing.T) {
		t.Parallel()
		t.Run("hostname by ip", func(t *testing.T) {
			t.Parallel()
			record := HostNameByIP(googleDNS)
			require.Equal(t, "dns.google.", record.Hostname)
		})
		t.Run("hostnames by ip", func(t *testing.T) {
			t.Parallel()
			hostnames := HostNamesByIP(googleDNS)
			require.Len(t, hostnames, 1)
		})
		t.Run("nameservers by hostname", func(t *testing.T) {
			t.Parallel()
			nameservers, err := NameServersByHostName("farishuskovic.dev")
			require.NoError(t, err)
			require.Len(t, nameservers, 4)
		})
		t.Run("network by hostname", func(t *testing.T) {
			t.Parallel()
			expected := net.ParseIP("8.0.0.0")
			record, err := NetworkByHost(googleDNS.String())
			require.NoError(t, err)
			require.Equal(t, expected.String(), record.NetworkIP.String())
		})
		t.Run("hostname and ip address using hostname", func(t *testing.T) {
			t.Parallel()
			expected, err := AddrByHostName("dns.google.")
			require.NoError(t, err)
			host, got, err := HostAndAddr("dns.google.")
			require.NoError(t, err)
			require.Equal(t, "dns.google.", host)
			require.Equal(t, expected.IP.String(), got.String())
		})
		t.Run("hostname and ip address using ip", func(t *testing.T) {
			t.Parallel()
			expected, err := AddrByHostName(googleDNS.String())
			require.NoError(t, err)
			host, got, err := HostAndAddr(googleDNS.String())
			require.NoError(t, err)
			require.Equal(t, host, "dns.google.")
			require.Equal(t, expected.IP.String(), got.String())
		})
		t.Run("internet service provider", func(t *testing.T) {
			t.Parallel()
			dec1st1992 := time.Date(1992, time.December, 1, 0, 0, 0, 0, time.UTC)
			_, network, err := net.ParseCIDR("8.8.8.0/24")
			require.NoError(t, err)
			expected := &InternetServiceProvider{
				Name:                    "GOOGLE, US",
				IP:                      &googleDNS,
				Country:                 "US",
				Registry:                "ARIN",
				IpRange:                 network,
				AutonomousServiceNumber: "AS15169",
				AllocatedAt:             &dec1st1992,
			}
			isp, err := ServiceProvider(&googleDNS)
			require.NoError(t, err)
			require.Equal(t, expected, isp)
		})
	})
	t.Run("ShouldFail", func(t *testing.T) {
		t.Parallel()
		t.Run("hostname by ip if invalid ip is input", func(t *testing.T) {
			t.Parallel()
			record := HostNameByIP(net.IP("invalid"))
			require.Empty(t, record.Hostname)
		})
		t.Run("hostnames by ip if invalid ip is input", func(t *testing.T) {
			t.Parallel()
			hostnames := HostNamesByIP(net.IP("invalid"))
			require.Nil(t, hostnames)
		})
		t.Run("addr by hostname if hostname is invalid", func(t *testing.T) {
			t.Parallel()
			addrs, err := AddrByHostName("invalid")
			require.Nil(t, addrs)
			require.Error(t, err)
		})
		t.Run("addr by hostname if hostname is invalid", func(t *testing.T) {
			t.Parallel()
			addrs, err := AddrByHostName("invalid")
			require.Nil(t, addrs)
			require.Error(t, err)
		})
	})
}
