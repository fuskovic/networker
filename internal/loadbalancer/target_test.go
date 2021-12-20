package loadbalancer

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTarget(t *testing.T) {
	for _, test := range []struct {
		name        string
		protocol    string
		host        string
		expectedErr error
	}{
		{
			name: "host with no port should fail",
			host: "hostwithnoport",
			expectedErr: fmt.Errorf(
				"expected %q to be formatted as host:port : %w", "hostwithnoport",
				errors.New("address hostwithnoport: missing port in address"),
			),
		},
		{
			name:        "unsupported protocol should fail",
			host:        "127.0.0.1:3000",
			protocol:    "ssh",
			expectedErr: ErrUnsupportedProtocol,
		},
		{
			name:        "valid target should pass",
			host:        "127.0.0.1:3000",
			protocol:    "http",
			expectedErr: nil,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			target, err := newTarget(
				&targetConfig{
					protocol: test.protocol,
					host:     test.host,
				},
			)
			if test.expectedErr != nil {
				require.Nil(t, target)
				require.Error(t, err)
				require.Equal(t, test.expectedErr.Error(), err.Error())
			} else {
				require.NotNil(t, target)
				require.NoError(t, err)
			}
		})
	}
}
