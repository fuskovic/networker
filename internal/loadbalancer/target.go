package loadbalancer

import (
	"crypto/tls"
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

func newTarget(protocol, host string) (*target, error) {
	// validate format
	if _, _, err := net.SplitHostPort(host); err != nil {
		return nil, fmt.Errorf("expected %q to be formatted as host:port : %w", host, err)
	}

	if protocol != "http" && protocol != "https" {
		return nil, ErrUnsupportedProtocol
	}

	url, err := url.Parse(fmt.Sprintf("%s://%s", protocol, host))
	if err != nil {
		return nil, fmt.Errorf("%q is an invalid url: %w", host, err)
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(url)
	if protocol == "https" {
		reverseProxy.Transport = &http.Transport{
			DialTLS: tlsDialer,
		}
	}

	return &target{
		ReverseProxy: reverseProxy,
		address:      host,
	}, nil
}

func tlsDialer(network, addr string) (net.Conn, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	tlsClient := tls.Client(conn,
		&tls.Config{ServerName: host},
	)

	if err := tlsClient.Handshake(); err != nil {
		conn.Close()
		return nil, err
	}

	state := tlsClient.ConnectionState()

	cert := state.PeerCertificates[0]
	if err := cert.VerifyHostname(host); err != nil {
		return nil, err
	}
	return tlsClient, nil
}
