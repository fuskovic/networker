package cmd

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
	unknown                 = "unknown"
	snapshotLen       int32 = 65535
	minSpacesSeq            = 12
	minSpacesProtocol       = 10
	minSpacesSrc            = 50
)

var (
	sep       = func(n int) string { return strings.Repeat(" ", n) }
	headerRow = fmt.Sprintf("SEQUENCE%sPROTOCOL%sSRC-MAC:SRC-IP:SRC-PORT%sDEST-MAC:DEST-IP:DEST-PORT", sep(5), sep(3), sep(28))
)

type (
	captureCmd struct {
		devices      []string
		seconds      int64
		outFile      string
		limit        bool
		numToCapture int64
	}

	row struct {
		srcIp, srcMac, srcPort, destIp, destMac, destPort, protocol string
		seq                                                         uint32
	}
)

// Spec returns a command spec containing a description of it's usage.
func (cmd *captureCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "capture",
		Usage: "[subcommand] [flags]",
		Desc:  "Capture network packets on specified devices.",
	}
}

// RegisterFlags initializes how a flag set is processed for a particular command.
func (cmd *captureCmd) RegisterFlags(fl *pflag.FlagSet) {
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
		flog.Error(err.Error())
		fl.Usage()
	}
}

func newRow() row {
	return row{
		srcIp:    unknown,
		srcMac:   unknown,
		srcPort:  unknown,
		destIp:   unknown,
		destMac:  unknown,
		destPort: unknown,
		protocol: unknown,
	}
}

func (r *row) format() string {
	pad := func(s string, n int) string {
		for len(s) < n {
			s += " "
		}
		return s
	}

	seq := pad(fmt.Sprintf("%d", r.seq), minSpacesSeq)
	src := pad(fmt.Sprintf("%s:%s:%s", r.srcMac, r.srcIp, r.srcPort), minSpacesSrc)
	protocol := pad(r.protocol, minSpacesProtocol)
	dst := fmt.Sprintf("%s:%s:%s", r.destMac, r.destIp, r.destPort)

	return fmt.Sprintf("%s %s %s %s", seq, protocol, src, dst)
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

	for _, d := range cmd.devices {
		for _, currentDevice := range allDevices {
			if currentDevice.Name == d {
				go cap(ctx, d, timeOut, packetChan)
			}
		}
	}

	flog.Info(headerRow)

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
			}

			unWrap(p)
			pktsCaptured++

			if limitReached(cmd.limit, cmd.numToCapture, pktsCaptured) {
				flog.Info("limit reached")
				cancel()
			}
		}
	}
	flog.Success("CAPTURE COMPLETE - CAPTURED %d PACKETS", pktsCaptured)
	return nil
}

func cap(ctx context.Context, device string, timeOut time.Duration, ch chan gopacket.Packet) error {
	handle, err := p.OpenLive(device, snapshotLen, false, timeOut)
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

func unWrap(packet gopacket.Packet) {
	row := newRow()

	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		ethernetPacket, decoded := ethernetLayer.(*layers.Ethernet)
		if !decoded {
			flog.Error("failed to decode ethernet layer")
		} else {
			row.srcMac = ethernetPacket.SrcMAC.String()
			row.destMac = ethernetPacket.DstMAC.String()
		}
	}

	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ipPacket, decoded := ipLayer.(*layers.IPv4)
		if !decoded {
			flog.Error("failed to decode IPv4 layer")
		} else {
			row.srcIp = ipPacket.SrcIP.String()
			row.destIp = ipPacket.DstIP.String()
			row.protocol = ipPacket.Protocol.String()
		}
	}

	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcpPacket, decoded := tcpLayer.(*layers.TCP)
		if !decoded {
			flog.Error("failed to decode TCP layer")
		} else {
			row.srcPort = tcpPacket.SrcPort.String()
			row.destPort = tcpPacket.DstPort.String()
			row.seq = tcpPacket.Seq
		}
	}

	if row.protocol != unknown {
		flog.Info(row.format())
	}

	// TODO : add flag so user is allowed to pass a network decryption key
	// so the payload can be decrypted.
	// applicationLayer := packet.ApplicationLayer()
	// if applicationLayer != nil {
	// 	fmt.Println("***Application layer***")
	// 	fmt.Printf("Payload: %s\n", applicationLayer.Payload())
	// }
}
