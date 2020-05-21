package cmd

import (
	"log"
	"net"

	"github.com/spf13/cobra"

	"github.com/fuskovic/networker/pkg/scan"
)

var (
	ip                         string
	ports                      []int
	upTo                       int
	tcpOnly, udpOnly, openOnly bool

	scanCmd = &cobra.Command{
		Use:     "scan",
		Aliases: []string{"s"},
		Example: scanExample,
		Short:   "Scan an IP for exposed ports.",
		Run: func(cmd *cobra.Command, args []string) {
			if net.ParseIP(ip) == nil {
				log.Printf("%s is not a valid IP address\n", ip)
				return
			}

			scanner := scan.NewScanner(ip, tcpOnly, udpOnly, openOnly)

			switch {
			case upTo > scan.TotalPorts:
				log.Printf("can not scan more than %d ports\n", scan.TotalPorts)
				return
			case upTo > scan.TotalPorts:
				log.Printf("can not scan more than %d ports\n", scan.TotalPorts)
				return
			case len(ports) > 0:
				scanner.ScanPorts(ports)
				log.Println("scan complete")
			case upTo > 0:
				scanner.ScanUpTo(upTo)
				log.Println("scan complete")
			default:
				scanner.ScanAllPorts()
				log.Println("scan complete")
			}
		},
	}
)

func init() {
	scanCmd.Flags().StringVar(&ip, "ip", "", "IP address to scan.")
	scanCmd.Flags().IntSliceVarP(&ports, "ports", "p", ports, "Specify a comma-separated list of ports to scan. (scans all ports if left unspecified)")
	scanCmd.Flags().IntVarP(&upTo, "up-to", "u", upTo, "Scan all ports up to a given port number.")
	scanCmd.Flags().BoolVarP(&tcpOnly, "tcp-only", "t", tcpOnly, "Only scan TCP ports.")
	scanCmd.Flags().BoolVar(&udpOnly, "udp-only", udpOnly, "Only scan UDP ports.")
	scanCmd.Flags().BoolVarP(&openOnly, "open-only", "o", openOnly, "Only print the ports that are open.")
	scanCmd.MarkFlagRequired("ip")
	Networker.AddCommand(scanCmd)
}
