package loadbalancer

import "errors"

var (
	// ErrNoPublicKey is returned when TLS is enabled but no public key is provided.
	ErrNoPublicKey = errors.New("no TLS public key provided")
	// ErrNoTrustCertificate is returned when TLS is enabled but no certificate is provided.
	ErrNoTrustCertificate = errors.New("no trust certificate provided")
	// ErrMinimumHostsUnmet is returned when the total amount of hosts provided is less than two.
	ErrMinimumHostsUnmet = errors.New("you must provide at least two hosts")
)

// Config collects the configuration parameters for the load balancer.
type Config struct {
	Hosts     []string
	Strategy  string
	EnableTLS bool
	PublicKey string
	Cert      string
}

func (cfg *Config) valid() error {
	if len(cfg.Hosts) < 2 {
		return ErrMinimumHostsUnmet
	}
	if cfg.EnableTLS && cfg.PublicKey == "" {
		return ErrNoPublicKey
	}
	if cfg.EnableTLS && cfg.Cert == "" {
		return ErrNoTrustCertificate
	}
	return nil
}
