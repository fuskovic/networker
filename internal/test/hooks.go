package test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

// WithNetworker is a pre-test hook that asserts the networker binary is globally installed before running the test.
// It's intended to be used in the cmd/networker pkg.
func WithNetworker(t *testing.T, name string, fn func(t *testing.T)) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		// make sure networker is installed
		cmd := exec.Command("which", "networker")
		require.NoError(t, cmd.Run(), "networker not installed")
		fn(t)
	})
}

// RequestFunc is the method signature the WithMockServer hook expects.
type RequestFunc func(t *testing.T, serverURL string)

// WithMockServer is a pre-test hook that starts the pre-configured mock server
// and passes it's URL to to the down stream test. It also handles teardown of the server
// and ensures the underlying test is parallelized.
func WithMockServer(t *testing.T, fn RequestFunc) {
	t.Helper()
	t.Parallel()
	testServer := newMockServer()
	defer testServer.Close()
	fn(t, testServer.URL)
}

type tlsFunc func(t *testing.T, s *TlsSuite)

// WithTlsSuite is pre-test hook for running TLS tests. The hook generates self-signed TLS certs before fn is ran
// and removes them after fn exists. The suite also provides utilities for initializing and spinning up TLS servers.
// WithTlsSuite parallelizes fn.
func WithTlsSuite(t *testing.T, name string, fn tlsFunc) {
	t.Helper()
	t.Parallel()
	t.Run(name, func(t *testing.T) {
		fn(t, newTlsTestSuite(t))
	})
}
