package capture

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	p "github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

const (
	snapshotLen int32 = 65535
	promiscuous       = false
)

// Packets captures network packets for the designated devices.
func Packets(designatedDevices []string, seconds, num int64, writer *pcapgo.Writer, isLimited bool) error {
	timeOut := time.Duration(seconds) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	allDevices, err := p.FindAllDevs()
	if err != nil {
		return fmt.Errorf("failed to find devices - err : %v", err)
	}

	packetChan := make(chan gopacket.Packet)

	for _, designatedDevice := range designatedDevices {
		for _, currentDevice := range allDevices {
			if currentDevice.Name == designatedDevice {
				go cap(ctx, currentDevice.Name, timeOut, packetChan)
			}
		}
	}

	var pktsCaptured int64

	for {
		select {
		case <-ctx.Done():
			fmt.Println("capture complete")
			return nil
		case p := <-packetChan:
			if writer != nil {
				if err := writer.WritePacket(p.Metadata().CaptureInfo, p.Data()); err != nil {
					return fmt.Errorf("failed to write to pcap - err : %v", err)
				}
			} else {
				fmt.Println(p)
			}

			pktsCaptured++
			if limitReached(isLimited, num, pktsCaptured) {
				fmt.Println("limit reached")
				cancel()
			}
		}
	}
}

func cap(ctx context.Context, device string, timeOut time.Duration, ch chan gopacket.Packet) error {
	handle, err := p.OpenLive(device, snapshotLen, promiscuous, timeOut)
	if err != nil {
		return err
	}
	defer handle.Close()

	time.Sleep(time.Second)

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

// NewWriter creates a new pcap and a writer for writing to it.
func NewWriter(outFile string) (*os.File, *pcapgo.Writer, error) {
	fileName := fmt.Sprintf("%s.pcap", outFile)
	f, err := os.Create(fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create %s", fileName)
	}

	w := pcapgo.NewWriter(f)
	if err := w.WriteFileHeader(uint32(snapshotLen), layers.LinkTypeEthernet); err != nil {
		return nil, nil, fmt.Errorf("failed to write pcap header - err : %v", err)
	}

	return f, w, nil
}

func limitReached(isLimited bool, limit, captured int64) bool {
	return isLimited && captured == limit
}
