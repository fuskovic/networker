package loadbalancer

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ErrUnsupportedProtocol is returned when http or https is not the designated protocol.
var ErrUnsupportedProtocol = errors.New("unsupported protocol")

type target struct {
	*httputil.ReverseProxy
	address  string
	hitCount int64
}

type targetConfig struct {
	protocol string
	host     string
	cert     []byte
	tlsCert  tls.Certificate
	isCA     bool
}

func newTarget(cfg *targetConfig) (*target, error) {
	// validate format
	host, _, err := net.SplitHostPort(cfg.host)
	if err != nil {
		return nil, fmt.Errorf("expected %q to be formatted as host:port : %w", cfg.host, err)
	}

	if cfg.protocol != "http" && cfg.protocol != "https" {
		return nil, ErrUnsupportedProtocol
	}

	url, err := url.Parse(fmt.Sprintf("%s://%s", cfg.protocol, cfg.host))
	if err != nil {
		return nil, fmt.Errorf("%q is an invalid url: %w", cfg.host, err)
	}

	tlsClientCfg := &tls.Config{ServerName: host}
	if cfg.isCA {
		certPool := x509.NewCertPool()
		if ok := certPool.AppendCertsFromPEM(cfg.cert); !ok {
			return nil, errors.New("failed to append cert to cert pool")
		}
		tlsClientCfg.RootCAs = certPool
	} else {
		tlsClientCfg.Certificates = append(tlsClientCfg.Certificates, cfg.tlsCert)
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(url)
	if cfg.protocol == "https" {
		reverseProxy.Transport = &http.Transport{
			DialTLS: func(network, addr string) (net.Conn, error) {
				host, _, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}

				conn, err := net.Dial(network, addr)
				if err != nil {
					return nil, err
				}

				tlsClient := tls.Client(conn, tlsClientCfg)
				if err := tlsClient.Handshake(); err != nil {
					conn.Close()
					return nil, fmt.Errorf("tls handshake failed: %w", err)
				}

				state := tlsClient.ConnectionState()
				cert := state.PeerCertificates[0]
				if err := cert.VerifyHostname(host); err != nil {
					return nil, fmt.Errorf("failed to verify hostname: %w", err)
				}
				return tlsClient, nil
			},
		}
	}

	return &target{
		ReverseProxy: reverseProxy,
		address:      cfg.host,
	}, nil
}
