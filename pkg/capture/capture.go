package capture

import (
	"fmt"
	"time"

	"github.com/google/gopacket"
	p "github.com/google/gopacket/pcap"
)

const (
	snapshotLen int32 = 1024
	promiscuous       = false
)

// Packets captures network packets for the designated devices.
func Packets(devices []string, seconds int64) error {
	timeOut := time.Duration(seconds) * time.Second

	allDevices, err := p.FindAllDevs()
	if err != nil {
		return err
	}

	if len(allDevices) == 0 {
		return fmt.Errorf("no connected devices")
	}

	for _, designatedDevice := range devices {
		for _, currentDevice := range allDevices {
			if currentDevice.Name == designatedDevice {
				go cap(currentDevice.Name, timeOut)
			}
		}
	}
	return nil
}

func cap(device string, timeOut time.Duration) error {
	handle, err := p.OpenLive(device, snapshotLen, promiscuous, timeOut)
	if err != nil {
		return err
	}
	defer handle.Close()

	src := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range src.Packets() {
		fmt.Println(packet)
	}
	return nil
}
