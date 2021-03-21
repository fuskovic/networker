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
	Kind     Kind   `table:"Kind"`
	Hostname string `table:"Hostname"`
	LocalIP  net.IP `table:"PrivateIP"`
	RemoteIP net.IP `table:"PublicIP"`
}
