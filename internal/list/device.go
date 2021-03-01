package list

import (
	"net"

	"cdr.dev/slog"
)

const (
	DeviceKindUnknown Kind = "unknown"
	DeviceKindRouter  Kind = "router"
	DeviceKindCurrent Kind = "you"
	DeviceKindPeer    Kind = "peer"
)

type Kind string

type Device struct {
	kind     Kind
	hostname string
	localIP  net.IP
	remoteIP net.IP
}

func (d *Device) Addr() string { return d.localIP.String() }

func (d *Device) Kind() Kind {
	switch d.kind {
	case DeviceKindCurrent, DeviceKindPeer, DeviceKindRouter:
		return d.kind
	}
	return DeviceKindUnknown
}

func (d *Device) Fields() []slog.Field {
	var fields []slog.Field
	if d.localIP != nil {
		fields = append(fields, slog.F("local-ip", d.localIP))
	}
	if d.hostname != "" {
		fields = append(fields, slog.F("hostname", d.hostname))
	}
	if d.remoteIP != nil {
		fields = append(fields, slog.F("remote-ip", d.remoteIP))
	}
	return fields
}
