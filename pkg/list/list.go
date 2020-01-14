package list

import (
	"fmt"

	p "github.com/google/gopacket/pcap"
)

// Device lists a device by its name.
func Device(name string) error {
	devices, err := p.FindAllDevs()
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
	devices, err := p.FindAllDevs()
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

func print(d p.Interface) {
	fmt.Printf("\nName: %s\nDescription: %s\n", d.Name, d.Description)
	for _, a := range d.Addresses {
		fmt.Printf("\n- IP address: %s\n- Subnet mask: %s\n", a.IP, a.Netmask)
	}
}

func match(a, b string) bool {
	return a == b
}
