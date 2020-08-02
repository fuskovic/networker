package scan

import "cdr.dev/slog"

type result struct {
	protocol string
	addr     string
	port     int
	open     bool
}

func (r *result) fields() []slog.Field {
	return []slog.Field{
		slog.F("port", r.port),
		slog.F("open", r.open),
	}
}
