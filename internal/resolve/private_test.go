package resolve

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrivate(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name            string
		ip              net.IP
		shouldBePrivate bool
	}{
		{
			name:            "loopbackIpV4",
			ip:              net.ParseIP("127.0.0.1"),
			shouldBePrivate: true,
		},
		{
			name:            "loopbackIpV6",
			ip:              net.ParseIP("::1"),
			shouldBePrivate: true,
		},
		{
			name:            "private class A",
			ip:              net.ParseIP("192.168.0.1"),
			shouldBePrivate: true,
		},
		{
			name:            "private class B",
			ip:              net.ParseIP("172.16.0.1"),
			shouldBePrivate: true,
		},
		{
			name:            "private class C",
			ip:              net.ParseIP("192.168.0.1"),
			shouldBePrivate: true,
		},
		{
			name:            "google DNS",
			ip:              net.ParseIP("8.8.8.8"),
			shouldBePrivate: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, test.shouldBePrivate, IsPrivate(&test.ip))
		})
	}
}
