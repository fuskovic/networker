package resolve

import (
	"net"
	"time"

	"cdr.dev/slog"
)

// InternetServiceProvider describes an internet service provider.
type InternetServiceProvider struct {
	name                    string
	ip                      *net.IP
	country                 string
	registry                string
	ipRange                 *net.IPNet
	autonomousServiceNumber string
	allocatedAt             *time.Time
}

// Fields returns a slice of formatted fields that can be used for logging.
func (isp *InternetServiceProvider) Fields() []slog.Field {
	var fields []slog.Field
	switch {
	case isp.name != "":
		fields = append(fields, slog.F("name", isp.name))
		fallthrough
	case isp.ip != nil:
		fields = append(fields, slog.F("ip-address", isp.ip))
		fallthrough
	case isp.ipRange != nil:
		fields = append(fields, slog.F("ip-range", isp.ipRange))
		fallthrough
	case isp.country != "":
		fields = append(fields, slog.F("country", isp.country))
		fallthrough
	case isp.registry != "":
		fields = append(fields, slog.F("registry", isp.registry))
		fallthrough
	case isp.autonomousServiceNumber != "":
		fields = append(fields, slog.F("autonomous-service-number", isp.autonomousServiceNumber))
		fallthrough
	case isp.allocatedAt != nil:
		fields = append(fields, slog.F("allocated-at", isp.allocatedAt))
	}
	return fields
}
