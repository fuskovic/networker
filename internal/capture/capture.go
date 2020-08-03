package capture

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	u "github.com/fuskovic/networker/internal/utils"
	pkt "github.com/google/gopacket"
	l "github.com/google/gopacket/layers"
	p "github.com/google/gopacket/pcap"
	pg "github.com/google/gopacket/pcapgo"
	"go.coder.com/flog"
)

const unknown = "unknown"

// Sniffer contains the fields that describe how to run the capture.
type Sniffer struct {
	Device  string
	Time    time.Duration
	File    string
	pktChan chan pkt.Packet
	Wide    bool
}

// Capture writes packets from the designated device to stdout and/or a pcap.
func (s *Sniffer) Capture() error {
	var captured int64

	ctx, cancel := context.WithTimeout(context.Background(), s.Time)
	defer cancel()

	s.pktChan = make(chan pkt.Packet)

	flog.Info("finding devices")

	devices, err := p.FindAllDevs()
	if err != nil {
		return fmt.Errorf("failed to find devices  : %v", err)
	}

	w, err := newWriter(s.File)
	if err != nil {
		return fmt.Errorf("failed to initialize new pcap writer : %v", err)
	}

	for _, d := range devices {
		if d.Name == s.Device {
			go s.sniff(ctx)
		}
	}

	log := sloghuman.Make(os.Stdout)

capture:
	for {
		select {
		case <-ctx.Done():
			break capture
		case p := <-s.pktChan:
			row := pktToRow(p, s.Wide)
			if row.Valid() {
				captured++
				log.Info(ctx, "pkt", row...)
			}

			md := p.Metadata().CaptureInfo

			if w != nil {
				err := w.WritePacket(md, p.Data())
				if err != nil {
					return fmt.Errorf("failed to write to pcap  : %v", err)
				}
			}
		}
	}
	log.Info(ctx, "capture complete", slog.F(
		"pkts-captured", captured),
	)
	return nil
}

func (s *Sniffer) sniff(ctx context.Context) error {
	h, err := p.OpenLive(s.Device,
		u.TotalPorts,
		false,
		p.BlockForever,
	)
	if err != nil {
		return err
	}
	defer h.Close()

	src := pkt.NewPacketSource(h, h.LinkType())
	flog.Info("capture started")

	for {
		select {
		case pkt := <-src.Packets():
			s.pktChan <- pkt
		case <-ctx.Done():
			return nil
		}
	}
}

func newWriter(fn string) (*pg.Writer, error) {
	if fn == "" {
		return nil, nil
	}

	if !strings.Contains(fn, ".pcap") {
		fn = fmt.Sprintf("%s.pcap", fn)
	}

	f, err := os.Create(fn)
	if err != nil {
		return nil, fmt.Errorf("failed to create %s", fn)
	}

	w := pg.NewWriter(f)
	err = w.WriteFileHeader(u.TotalPorts, l.LinkTypeEthernet)
	return w, err
}

func pktToRow(p pkt.Packet, wide bool) u.Row {
	var (
		r       u.Row
		proto   = unknown
		srcMac  = unknown
		dstMac  = unknown
		srcIP   = unknown
		dstIP   = unknown
		srcPort = unknown
		dstPort = unknown
		seq     = -1
	)

	ethLayer := p.Layer(l.LayerTypeEthernet)
	if ethLayer != nil {
		ethPkt, ok := ethLayer.(*l.Ethernet)
		if ok {
			srcMac = ethPkt.SrcMAC.String()
			dstMac = ethPkt.DstMAC.String()
		}
	}

	if wide {
		r.Add("src-mac", srcMac)
		r.Add("dst-mac", dstMac)
	}

	ipLayer := p.Layer(l.LayerTypeIPv4)
	if ipLayer != nil {
		ipPkt, ok := ipLayer.(*l.IPv4)
		if ok {
			srcIP = ipPkt.SrcIP.String()
			dstIP = ipPkt.DstIP.String()
			proto = ipPkt.Protocol.String()
		}
	}

	r.Add("src-ip", srcIP)
	r.Add("dst-ip", dstIP)
	if wide {
		r.Add("src-host", u.HostNameByIP(srcIP))
		r.Add("dst-host", u.HostNameByIP(dstIP))
	}

	tcpLayer := p.Layer(l.LayerTypeTCP)
	if tcpLayer != nil {
		tcpPkt, ok := tcpLayer.(*l.TCP)
		if ok {
			seq = int(tcpPkt.Seq)
			srcPort = tcpPkt.SrcPort.String()
			dstPort = tcpPkt.DstPort.String()
		}
	}

	r.Add("src-port", srcPort)
	r.Add("dst-port", dstPort)
	r.Add("proto", proto)
	if wide {
		r.Add("seq", seq)
	}
	return r
}