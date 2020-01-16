package capture

import (
	"context"
	"fmt"
	"time"

	"github.com/google/gopacket"
	p "github.com/google/gopacket/pcap"
)

const (
	snapshotLen int32 = 65535
	promiscuous       = false
)

// Packets captures network packets for the designated devices.
func Packets(designatedDevices []string, seconds int64) error {
	timeOut := time.Duration(seconds) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	allDevices, err := p.FindAllDevs()
	if err != nil {
		return fmt.Errorf("failed to find devices - err : %v", err)
	}

	if len(allDevices) == 0 {
		return fmt.Errorf("no connected devices")
	}

	if len(designatedDevices) == 0 {
		return fmt.Errorf("no designated devices")
	}

	packetChan := make(chan gopacket.Packet)

	go func(ctx context.Context) {
		for _, designatedDevice := range designatedDevices {
			for _, currentDevice := range allDevices {
				if currentDevice.Name == designatedDevice {
					go cap(ctx, currentDevice.Name, timeOut, packetChan)
				}
			}
		}
	}(ctx)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("capture complete")
			return nil
		case packet := <-packetChan:
			fmt.Println(packet)
		}
	}
}

func cap(ctx context.Context, device string, timeOut time.Duration, ch chan gopacket.Packet) error {
	handle, err := p.OpenLive(device, snapshotLen, promiscuous, timeOut)
	if err != nil {
		return err
	}
	defer handle.Close()

	src := gopacket.NewPacketSource(handle, handle.LinkType())

	for {
		select {
		case packet := <-src.Packets():
			ch <- packet
		case <-ctx.Done():
			return nil
		}
	}
}
