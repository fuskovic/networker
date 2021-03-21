package ports

import "net"

type Scan struct {
	IpAddress      net.IP `table:"IP-Address"`
	Hostname       string `table:"Hostname"`
	FoundOpenPorts bool   `table:"Found-Open-Ports"`
	OpenPorts      []int  `table:"Open-Ports"`
}
