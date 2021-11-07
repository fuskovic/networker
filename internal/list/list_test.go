package list

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListDevices(t *testing.T) {
	t.Parallel()
	// We only check the current device exists in the list of devices
	// because its the only device guaranteed to be on the LAN.
	// Asserting against a static list of devices is not reliable because those
	// devices may not be on the network at anymore at any given time.
	t.Run("find current device in list of devices", func(t *testing.T) {
		ctx := context.Background()
		currentDevice, err := getCurrentDevice(ctx)
		require.NoError(t, err)
		devices, err := Devices(ctx)
		require.NoError(t, err)
		var foundDevice bool
		for _, d := range devices {
			if d.LocalIP.String() == currentDevice.LocalIP.String() {
				foundDevice = true
			}
		}
		require.True(t, foundDevice)
	})
}
