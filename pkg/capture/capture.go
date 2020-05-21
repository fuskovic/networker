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

// Config collects the command parameters for the capture sub-command.
type Config struct {
	Devices        []string
	Seconds        int64
	OutFile        string
	Limit, Verbose bool
	NumToCapture   int64
}

// Run executes the command logic for the capture package.
func Run(cfg *Config) error {
	if len(cfg.Devices) == 0 {
		return fmt.Errorf("no designated devices")
	}

	if cfg.Seconds < 5 {
		return fmt.Errorf("capture must be at least 5 seconds long - your input : %d", cfg.Seconds)
	}

	if cfg.Limit && cfg.NumToCapture < 1 {
		return fmt.Errorf("use of --limit flag without use of --num flag\nPlease specify number of packets to limit capture\nminimum is 1")
	}

	if err := start(cfg); err != nil {
		return fmt.Errorf("error during packet capture : %v", err)
	}
	return nil
}

func start(cfg *Config) error {
	var writer *pcapgo.Writer
	timeOut := time.Duration(cfg.Seconds) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	allDevices, err := p.FindAllDevs()
	if err != nil {
		return fmt.Errorf("failed to find devices - err : %v", err)
	}

	packetChan := make(chan gopacket.Packet)
	var pktsCaptured int64

	if cfg.OutFile != "" {
		file, w, err := newWriter(cfg.OutFile)
		if err != nil {
			return fmt.Errorf("failed to create a new writer - err : %v", err)
		}
		defer file.Close()
		writer = w
	}

	if cfg.Verbose {
		go logProgress(ctx, timeOut, cfg.Limit, &pktsCaptured, cfg.NumToCapture)
	}

	for _, d := range cfg.Devices {
		for _, currentDevice := range allDevices {
			if currentDevice.Name == d {
				go cap(ctx, currentDevice.Name, timeOut, packetChan)
			}
		}
	}

capture:
	for {
		select {
		case <-ctx.Done():
			break capture
		case p := <-packetChan:
			if writer != nil {
				if err := writer.WritePacket(p.Metadata().CaptureInfo, p.Data()); err != nil {
					return fmt.Errorf("failed to write to pcap - err : %v", err)
				}
			} else {
				// TODO : clean up this output
				fmt.Println(p)
			}

			pktsCaptured++
			if limitReached(cfg.Limit, cfg.NumToCapture, pktsCaptured) {
				fmt.Println("limit reached")
				cancel()
			}
		}
	}
	return nil
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

func newWriter(outFile string) (*os.File, *pcapgo.Writer, error) {
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

func logProgress(ctx context.Context, d time.Duration, isLimited bool, captured *int64, limit int64) {
	start := time.Now()
	end := start.Add(d)

	for start.Unix() < end.Unix() {
		elapsed := time.Since(start).Truncate(time.Second)
		time.Sleep(500 * time.Millisecond)
		output := fmt.Sprintf("\r%v/%v elapsed", elapsed, d)
		if isLimited {
			output += fmt.Sprintf(" - %d/%d captured", *captured, limit)
		}
		fmt.Print(output)
	}
}
