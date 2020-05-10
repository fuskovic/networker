package cmd

import (
	"log"

	"github.com/fuskovic/github.com/networker/pkg/scan"
	"github.com/spf13/cobra"
)

var (
	host             string
	ports            []int
	upTo             int
	tcpOnly, udpOnly bool
	scanEx           = "TODO : example scan command"

	scanCmd = &cobra.Command{
		Use:     "scan",
		Aliases: []string{"s"},
		Example: scanEx,
		Short:   "scan for exposed ports on a designated IP",
		Run: func(cmd *cobra.Command, args []string) {
			scanner := scan.NewScanner(host, tcpOnly, udpOnly)

			if len(ports) > 0 {
				scanner.ScanPorts(ports, tcpOnly, udpOnly)
				return
			}

			if upTo != 0 {
				scanner.ScanUpTo(upTo, tcpOnly, udpOnly)
				return
			}

			log.Println("scanning all ports...")
			scanner.ScanAllPorts(host, tcpOnly, udpOnly)
		},
	}
)

func init() {
	scanCmd.Flags().StringVar(&host, "host", "", "IP address of host to scan")
	scanCmd.Flags().IntSliceVar(&ports, "ports", ports, "explicitly specify which ports you want scanned (comma separated). If not specified, all ports will be scanned")
	scanCmd.Flags().IntVarP(&upTo, "upto", "u", upTo, "scan all ports up to a specified value")
	scanCmd.Flags().BoolVar(&tcpOnly, "tcp-only", tcpOnly, "enable to scan only tcp ports")
	scanCmd.Flags().BoolVar(&udpOnly, "udp-only", udpOnly, "enable to scan only udp ports")
	scanCmd.MarkFlagRequired("host")
	Networker.AddCommand(scanCmd)
}
