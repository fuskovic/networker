package list

import (
	"fmt"
	"net"
)

// LocalIP lists the local IP address of the node executing this command.
func LocalIP() {
	// TODO: implement
}

// RemoteIP lists the remote IP address of the node executing this command.
func RemoteIP() {
	// TODO: implement
}

// Router lists the IP address of the gateway on this subnet.
func Router() {
	// TODO: implement
}

// Device lists a device by its name.
func Device(name string) error {
	devices, err := net.Interfaces()
	if err != nil {
		return err
	}

	lastDevice := len(devices)

	if lastDevice == 0 {
		return fmt.Errorf("no devices found")
	}

	for i, d := range devices {
		if d.Name == name {
			print(d)
			return nil
		}

		if i+1 == lastDevice && !match(d.Name, name) {
			return fmt.Errorf("device : %s not found", name)
		}
	}

	return nil
}

// AllDevices lists all connected network interfaces.
func AllDevices() error {
	devices, err := net.Interfaces()
	if err != nil {
		return err
	}

	if len(devices) == 0 {
		return fmt.Errorf("no devices found")
	}

	fmt.Printf("found %d devices", len(devices))
	for _, d := range devices {
		print(d)
	}
	return nil
}

func print(d net.Interface) {
	fmt.Printf("\nIndex : %d\nName: %s\nHardware Address: %s\nMTU : %d\nFlags : %s\n",
		d.Index,
		d.Name,
		d.HardwareAddr.String(),
		d.MTU,
		d.Flags.String(),
	)
	addrs, _ := d.Addrs()
	for _, a := range addrs {
		fmt.Printf("\n- IP address: %s\n- Network: %s\n", a.String(), a.Network())
	}
	mcAddrs, _ := d.MulticastAddrs()
	for _, ma := range mcAddrs {
		fmt.Printf("\n- IP address: %s\n- Network: %s\n", ma.String(), ma.Network())
	}
}

func match(a, b string) bool {
	return a == b
}
