package test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// ProjectRoot is a utility function for returning root path of the networker source code on this machine.
// It's primary purpose is to help construct absolute paths needed by tests.
func ProjectRoot(t *testing.T) string {
	t.Helper()
	output, err := exec.Command("git", "rev-parse", "--show-toplevel").CombinedOutput()
	require.NoError(t, err)
	return strings.Replace(string(output), "\n", "", 1)
}
