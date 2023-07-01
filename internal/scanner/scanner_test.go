package scanner

import (
	"context"
	"net"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScanner(t *testing.T) {
	t.Skip("flakey")
	// start a listener on any available local port
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer l.Close()

	// check which port the listener is running to determine how many ports we need to scan
	host, portStr, err := net.SplitHostPort(l.Addr().String())
	require.NoError(t, err)
	port, err := strconv.Atoi(portStr)
	require.NoError(t, err)
	var shouldScanAll bool
	if port > wellKnownPorts {
		shouldScanAll = true
	}

	// initialize a new scanner and scan localhost
	hostsToScan := []string{host}
	results, err := New(hostsToScan, shouldScanAll).Scan(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, len(results))

	// assert that the port the listener was started on is listed as an open port in the results
	var foundPort bool
	for _, openPort := range results[0].Ports {
		if openPort == port {
			foundPort = true
		}
	}
	require.True(t, foundPort)
}
