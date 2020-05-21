package cmd

import (
	"log"

	"github.com/fuskovic/networker/pkg/lookup"
	"github.com/spf13/cobra"
)

type lookUpFunc func(string) error

var (
	hostName, ipAddress, nameServer, network string

	supportedLookUps = map[*string]lookUpFunc{
		&hostName:   lookup.HostNamesByIP,
		&ipAddress:  lookup.IPsByHostName,
		&nameServer: lookup.NameServersByHostName,
		&network:    lookup.NetworkByHostName,
	}

	lookUpCmd = &cobra.Command{
		Use:     "lookup",
		Aliases: []string{"lu"},
		Example: lookUpExample,
		Short:   "Lookup hostnames, IP addresses, nameservers, and general network information.",
		Run: func(cmd *cobra.Command, args []string) {
			for value, lookUp := range supportedLookUps {
				if *value != "" {
					if err := lookUp(*value); err != nil {
						log.Printf("failed lookup\nerror : %v\n", err)
						cmd.Usage()
					}
				}
			}
		},
	}
)

func init() {
	lookUpCmd.Flags().StringVarP(&network, "network", "n", "", "Look up the network a given hostname belongs to.")
	lookUpCmd.Flags().StringVarP(&ipAddress, "addresses", "a", "", "Look up IP addresses for a given hostname.")
	lookUpCmd.Flags().StringVarP(&nameServer, "nameservers", "s", "", "Look up nameservers for a given hostname.")
	lookUpCmd.Flags().StringVar(&hostName, "hostnames", "", "Look up hostnames for a given IP address.")
	Networker.AddCommand(lookUpCmd)
}
