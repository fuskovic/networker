package list

import "net"

const (
	DeviceKindUnknown Kind = "unknown"
	DeviceKindRouter  Kind = "router"
	DeviceKindCurrent Kind = "current-device"
	DeviceKindPeer    Kind = "peer"
)

type Kind string

type Device struct {
	Kind     Kind   `json:"kind" table:"Kind"`
	Hostname string `json:"hostname" table:"Hostname"`
	LocalIP  net.IP `json:"local_ip" table:"LocalIP"`
	RemoteIP net.IP `json:"remote_ip,omitempty" table:"RemoteIP"`
}
