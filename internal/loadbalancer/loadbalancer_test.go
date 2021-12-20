package loadbalancer

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fuskovic/networker/internal/test"
	"github.com/stretchr/testify/require"
)

func TestLoadBalancer(t *testing.T) {
	t.Run("ShouldFail", func(t *testing.T) {
		t.Run("invalid config should fail", func(t *testing.T) {
			_, err := New(&Config{Hosts: []string{"127.0.0.1"}})
			require.Error(t, err)
			require.Equal(t, err,
				fmt.Errorf("invalid load balancer configuration: %w",
					ErrMinimumHostsUnmet,
				))
		})
	})
	t.Run("ShouldPass", func(t *testing.T) {
		test.WithTlsSuite(t, "round robin", func(t *testing.T, suite *test.TlsSuite) {
			// start three mock tls servers
			tlsServers, teardown := suite.NewStartedMockTlsServers(t, 3)
			defer teardown()

			// setup the load balaner configuration
			cfg := &Config{
				Cert:      suite.CACert,
				EnableTLS: true,
				IsTest:    true,
			}

			for _, s := range tlsServers {
				cfg.Hosts = append(cfg.Hosts, s.URL())
			}

			// initialize the load balancer and start it
			lb, err := New(cfg)
			require.NoError(t, err)
			tlsServer := suite.NewUnstartedMockTlsServer(t)
			tlsServer.SetHandler(lb)
			tlsServer.StartTLS()
			defer tlsServer.Close()

			// initialize the client
			tlsClient := http.Client{
				Transport: &http.Transport{
					TLSClientConfig: suite.ClientConfig,
				},
			}
			defer tlsClient.CloseIdleConnections()

			//make 9 http requests to the load balancer
			for i := 0; i < 9; i++ {
				w := httptest.NewRecorder()
				w.Result()
				r := httptest.NewRequest(http.MethodGet, tlsServer.URL(), nil)
				r.RequestURI = ""
				resp, err := tlsClient.Do(r)
				require.NoError(t, err)
				require.True(t, resp.TLS.HandshakeComplete)
				require.Equal(t, http.StatusOK, resp.StatusCode)
			}
			// assert that each server had a hit count of 3
			for _, s := range tlsServers {
				require.Equal(t, s.HitCount(), 3)
			}
		})
	})
	// t.Run("valid non-tls enabled config should pass", func(t *testing.T) {
	// 	lb, err := New(
	// 		&Config{
	// 			Hosts:     []string{},
	// 			Strategy:  "",
	// 			EnableTLS: false,
	// 			PublicKey: "",
	// 			Cert:      "",
	// 		},
	// 	)
	// 	require.NoError(t, err)
	// })
}
