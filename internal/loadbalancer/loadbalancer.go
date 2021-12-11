package loadbalancer

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	ErrNoPublicKey         = errors.New("no TLS public key provided")
	ErrMinimumTargetsUnmet = errors.New("targets provided less than minimum of 2")
)

type (
	Balancer interface {
		http.Handler
		Balance() chan error
	}

	Config struct {
		Targets   []string
		Strategy  string
		EnableTLS bool
		PublicKey string
		Port      string
	}

	loadBalancer struct {
		port    string
		targets []*target
	}

	target struct {
		*httputil.ReverseProxy
		address  string
		hitCount int64
	}
)

func New(cfg *Config) (Balancer, error) {
	if len(cfg.Targets) < 2 {
		return nil, ErrMinimumTargetsUnmet
	}

	if cfg.EnableTLS && cfg.PublicKey == "" {
		return nil, ErrNoPublicKey
	}

	lb := &loadBalancer{port: cfg.Port}

	for i := range cfg.Targets {
		host, port, err := net.SplitHostPort(cfg.Targets[i])
		if err != nil {
			return nil, fmt.Errorf("%q is an invalid target address: %w", cfg.Targets[i], err)
		}

		protocol := "http"
		if cfg.EnableTLS {
			protocol = "https"
		}

		addr := net.JoinHostPort(host, port)
		lb.targets = append(lb.targets,
			&target{
				ReverseProxy: httputil.NewSingleHostReverseProxy(
					&url.URL{
						Scheme: protocol,
						Host:   addr,
					},
				),
				address: addr,
			},
		)

		// if cfg.EnableTLS {
		// 	lb.targets[host].Transport = &httlb.Transport{DialTLS: tlsDialer}
		// }
	}
	return lb, nil
}

func (lb *loadBalancer) Balance() chan error {
	c := make(chan error, 1)
	go func() {
		c <- http.ListenAndServe(lb.port, lb)
	}()
	return c
}

func (lb *loadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newTarget := lb.targets[0]
	for _, target := range lb.targets {
		if target.hitCount < newTarget.hitCount {
			newTarget = target
		}
	}

	newTarget.hitCount++
	log.Printf("from: %q to: %q hit-count: %d",
		r.RemoteAddr,
		newTarget.address,
		newTarget.hitCount,
	)
	newTarget.ServeHTTP(w, r)
}
