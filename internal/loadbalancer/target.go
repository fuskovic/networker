package loadbalancer

import (
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
		// TODO: make sure transport supports TLS
		reverseProxy.Transport = &http.Transport{}
	}

	return &target{
		ReverseProxy: reverseProxy,
		address:      host,
	}, nil
}
