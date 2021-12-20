package loadbalancer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	for _, test := range []struct {
		name        string
		config      Config
		expectedErr error
	}{
		{
			name:        "config with less than 2 target hosts should fail",
			config:      Config{Hosts: []string{"127.0.0.1"}},
			expectedErr: ErrMinimumHostsUnmet,
		},
		{
			name: "valid tls cert should pass",
			config: Config{
				Hosts: []string{
					"127.0.0.1",
					"doesntmatter",
				},
			},
			expectedErr: nil,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := test.config.valid()
			require.Equal(t, err, test.expectedErr)
		})
	}
}
