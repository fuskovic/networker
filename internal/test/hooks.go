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
