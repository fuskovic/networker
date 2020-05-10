package cmd

import "github.com/spf13/cobra"

var (
	host   string
	ports  []int
	scanEx = "TODO : example scan command"

	scanCmd = &cobra.Command{
		Use:     "scan",
		Aliases: []string{"s"},
		Example: scanEx,
		Short:   "scan for exposed ports on a designated IP",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO : scan specified ports for host
			// TODO : scan all ports of host if ports unspecifie
		},
	}
)

func init() {
	scanCmd.Flags().StringVar(&host, "host", "", "IP address of host to scan")
	scanCmd.Flags().IntSliceVar(&ports, "ports", ports, "explicitly specify which ports you want scanned (comma separated). If not specified, all ports will be scanned")
	scanCmd.MarkFlagRequired("host")
	Networker.AddCommand(scanCmd)
}
