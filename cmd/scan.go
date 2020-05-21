package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fuskovic/networker/pkg/scan"
)

var (
	scanCfg = &scan.Config{}

	scanCmd = &cobra.Command{
		Use:     "scan",
		Aliases: []string{"s"},
		Example: scanExample,
		Short:   "Scan an IP for exposed ports.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := scan.Run(scanCfg); err != nil {
				fmt.Println(err)
				cmd.Usage()
			}
		},
	}
)

func init() {
	scanCmd.Flags().StringVar(&scanCfg.IP, "ip", "", "IP address to scan.")
	scanCmd.Flags().IntSliceVarP(&scanCfg.Ports, "ports", "p", scanCfg.Ports, "Specify a comma-separated list of ports to scan. (scans all ports if left unspecified)")
	scanCmd.Flags().IntVarP(&scanCfg.UpTo, "up-to", "u", scanCfg.UpTo, "Scan all ports up to a given port number.")
	scanCmd.Flags().BoolVarP(&scanCfg.TCPonly, "tcp-only", "t", scanCfg.TCPonly, "Only scan TCP ports.")
	scanCmd.Flags().BoolVar(&scanCfg.UDPonly, "udp-only", scanCfg.UDPonly, "Only scan UDP ports.")
	scanCmd.Flags().BoolVarP(&scanCfg.OpenOnly, "open-only", "o", scanCfg.OpenOnly, "Only print the ports that are open.")
	scanCmd.MarkFlagRequired("ip")
	Networker.AddCommand(scanCmd)
}
