package capture

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	p "github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"
)

const (
	snapshotLen int32 = 65535
	promiscuous       = false
)

var stars = strings.Repeat("*", 30)

type captureCmd struct {
	devices        []string
	seconds        int64
	outFile        string
	limit, verbose bool
	numToCapture   int64
}

// Spec returns a command spec containing a description of it's usage.
func (cmd *captureCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "capture",
		Usage: "TODO: ADD USAGE",
		Desc:  "Capture network packets on specified devices.",
	}
}

// RegisterFlags initializes how a flag set is processed for a particular command.
func (cmd *captureCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.BoolVarP(&cmd.verbose, "verbose", "v", cmd.verbose, "Enable verbose logging.")
	fl.Int64VarP(&cmd.seconds, "seconds", "s", cmd.seconds, "Amount of seconds to run capture for.")
	fl.StringSliceVarP(&cmd.devices, "devices", "d", cmd.devices, "Comma-separated list of devices to capture packets on.")
	fl.StringVarP(&cmd.outFile, "out", "o", cmd.outFile, "Name of an output file to write the packets to.")
	fl.BoolVarP(&cmd.limit, "limit", "l", cmd.limit, "Limit the number of packets to capture. (must be used with the --num flag)")
	fl.Int64VarP(&cmd.numToCapture, "num", "n", cmd.numToCapture, "Number of total packets to capture across all devices.")
}

// Run validates the flagset and if successful, runs the packet capture session accordingly.
func (cmd *captureCmd) Run(fl *pflag.FlagSet) {
	var err error

	switch {
	case len(cmd.devices) == 0:
		err = fmt.Errorf("no designated devices")
	case cmd.seconds < 5:
		err = fmt.Errorf("capture must be at least 5 seconds long - your input : %d", cmd.seconds)
	case cmd.limit && cmd.numToCapture < 1:
		err = fmt.Errorf("use of --limit flag without use of --num flag\nPlease specify number of packets to limit capture\nminimum is 1")
	default:
		err = start(cmd)
	}

	if err != nil {
		flog.Fatal(err.Error())
	}
}

func start(cmd *captureCmd) error {
	var writer *pcapgo.Writer
	timeOut := time.Duration(cmd.seconds) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	allDevices, err := p.FindAllDevs()
	if err != nil {
		return fmt.Errorf("failed to find devices - err : %v", err)
	}

	packetChan := make(chan gopacket.Packet)
	var pktsCaptured int64

	if cmd.outFile != "" {
		file, w, err := newWriter(cmd.outFile)
		if err != nil {
			return fmt.Errorf("failed to create a new writer - err : %v", err)
		}
		defer file.Close()
		writer = w
	}

	if cmd.verbose {
		go logProgress(ctx, timeOut, cmd.limit, &pktsCaptured, cmd.numToCapture)
	}

	for _, d := range cmd.devices {
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
				fmt.Println(stars)
				unWrap(p)
			}

			pktsCaptured++
			if limitReached(cmd.limit, cmd.numToCapture, pktsCaptured) {
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

func unWrap(packet gopacket.Packet) {
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		fmt.Println("***Ethernet layer***")
		ethernetPacket, decoded := ethernetLayer.(*layers.Ethernet)
		if !decoded {
			fmt.Println("failed to decode ethernet layer")
		} else {
			fmt.Println("Source MAC: ", ethernetPacket.SrcMAC)
			fmt.Println("Destination MAC: ", ethernetPacket.DstMAC)
			fmt.Println("Ethernet type: ", ethernetPacket.EthernetType)
		}
	}

	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		fmt.Println("***IPv4 layer***")
		ipPacket, decoded := ipLayer.(*layers.IPv4)
		if !decoded {
			fmt.Println("failed to decode IPv4 layer")
		} else {
			fmt.Printf("From %s to %s\n", ipPacket.SrcIP, ipPacket.DstIP)
			fmt.Println("Protocol: ", ipPacket.Protocol)
		}
	}

	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		fmt.Println("***TCP layer***")
		tcpPacket, decoded := tcpLayer.(*layers.TCP)
		if !decoded {
			fmt.Println("failed to decode TCP layer")
		} else {
			fmt.Printf("From port %d to %d\n", tcpPacket.SrcPort, tcpPacket.DstPort)
			fmt.Println("Sequence number: ", tcpPacket.Seq)
		}
	}

	// TODO : add flag so user is allowed to pass a network decryption key
	// so the payload can be decrypted.
	// applicationLayer := packet.ApplicationLayer()
	// if applicationLayer != nil {
	// 	fmt.Println("***Application layer***")
	// 	fmt.Printf("Payload: %s\n", applicationLayer.Payload())
	// }
}
