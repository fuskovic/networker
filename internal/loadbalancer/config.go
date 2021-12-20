package loadbalancer

import (
	"errors"
	"strings"
)

// ErrMinimumHostsUnmet is returned when the total amount of hosts provided is less than two.
var ErrMinimumHostsUnmet = errors.New("you must provide at least two hosts")

// Config collects the configuration parameters for the load balancer.
type Config struct {
	Hosts     []string
	Strategy  string
	EnableTLS bool
	IsTest    bool
	Cert      []byte
}

func (cfg *Config) valid() error {
	if len(cfg.Hosts) < 2 {
		return ErrMinimumHostsUnmet
	}
	// normalize the hosts by removing the provided protocols
	for i, host := range cfg.Hosts {
		// remove the protocol
		cfg.Hosts[i] = strings.ReplaceAll(host, "http://", "")
		cfg.Hosts[i] = strings.ReplaceAll(host, "https://", "")
	}
	return nil
}
