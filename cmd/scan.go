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
	scanEx                     = "TODO : example scan command"

	scanCmd = &cobra.Command{
		Use:     "scan",
		Aliases: []string{"s"},
		Example: scanEx,
		Short:   "scan for exposed ports on a designated IP",
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
	scanCmd.Flags().IntSliceVar(&ports, "ports", ports, "explicitly specify which ports you want scanned (comma separated). If not specified, all ports will be scanned unless up-to flag is specified.")
	scanCmd.Flags().IntVarP(&upTo, "up-to", "u", upTo, "scan all ports up to a specified value")
	scanCmd.Flags().BoolVar(&tcpOnly, "tcp-only", tcpOnly, "enable to scan only tcp ports")
	scanCmd.Flags().BoolVar(&udpOnly, "udp-only", udpOnly, "enable to scan only udp ports")
	scanCmd.Flags().BoolVar(&openOnly, "open-only", openOnly, "enable to only log open ports")
	scanCmd.MarkFlagRequired("ip")
	Networker.AddCommand(scanCmd)
}
