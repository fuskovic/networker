package loadbalancer

import (
	"fmt"
	"log"
	"net/http"
)

type loadBalancer struct {
	targets []*target
}

// New initializes a new handler that can be used to load balance
// network traffic across the hosts designated in the config.
func New(cfg *Config) (http.Handler, error) {
	if err := cfg.valid(); err != nil {
		return nil, fmt.Errorf("invalid load balancer configuration: %w", err)
	}

	lb := new(loadBalancer)
	protocol := "http"
	if cfg.EnableTLS {
		protocol = "https"
	}

	for _, host := range cfg.Hosts {
		target, err := newTarget(protocol, host)
		if err != nil {
			return nil, fmt.Errorf("failed to mount target host %q to load balancer: %w", host, err)
		}
		lb.targets = append(lb.targets, target)
	}
	return lb, nil
}

// ServeHTTP routes requests to the target host with the lowest hit-count.
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
