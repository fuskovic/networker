package cmd

import (
	"fmt"
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
	specificPortsExample       = "scan only a specified set of TCP ports and only log if they're open:\nnetworker scan --ip <someIPaddress> --ports 22,80,3389 --open-only\n"
	longScanExample            = "scan all TCP ports up to port 1024 and only log status if they're open:\nnetworker scan --ip <someIPaddress> --up-to 1024 --tcp-only --open-only\n"
	shortScanExample           = "\nshort form: networker s --ip <someIPaddress> --up-to 1024 -t -o\n"
	scanExample                = fmt.Sprintf("%s\n%s\n%s\n", specificPortsExample, longScanExample, shortScanExample)

	scanCmd = &cobra.Command{
		Use:     "scan",
		Aliases: []string{"s"},
		Example: scanExample,
		Short:   "scan for exposed ports on a designated IP.",
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
	scanCmd.Flags().StringVar(&ip, "ip", "", "IP address to scan")
	scanCmd.Flags().IntSliceVarP(&ports, "ports", "p", ports, "explicitly specify which ports you want scanned (comma separated). If not specified, all ports will be scanned unless up-to flag is specified.")
	scanCmd.Flags().IntVarP(&upTo, "up-to", "u", upTo, "scan all ports up to a specified value")
	scanCmd.Flags().BoolVarP(&tcpOnly, "tcp-only", "t", tcpOnly, "enable to scan only tcp ports")
	scanCmd.Flags().BoolVar(&udpOnly, "udp-only", udpOnly, "enable to scan only udp ports")
	scanCmd.Flags().BoolVarP(&openOnly, "open-only", "o", openOnly, "enable to only log open ports")
	scanCmd.MarkFlagRequired("ip")
	Networker.AddCommand(scanCmd)
}
