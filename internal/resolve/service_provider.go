package resolve

import (
	"net"
	"time"
)

// InternetServiceProvider describes an internet service provider.
type InternetServiceProvider struct {
	Name                    string     `json:"name" table:"Name"`
	IP                      *net.IP    `json:"ip_address" table:"IP"`
	Country                 string     `json:"country" table:"Country"`
	Registry                string     `json:"registry" table:"Registry"`
	IpRange                 *net.IPNet `json:"ip_range" table:"IP-Range"`
	AutonomousServiceNumber string     `json:"autonomous_service_number" table:"ASN"`
	AllocatedAt             *time.Time `json:"allocated_at" table:"AllocatedAt"`
}
