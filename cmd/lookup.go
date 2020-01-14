package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/fuskovic/networker/pkg/lookup"
)

type lookUpFunc func(string) error

var (
	hostName, ipAddress, nameServer, domain string

	supportedLookUps = map[*string]lookUpFunc{
		&hostName : lookup.HostNamesByIP,
		&ipAddress: lookup.IPsByHostName,
		&nameServer:lookup.NameServersByHostName,
		&domain:	lookup.MxRecordsForDomain,
	}

	lookUpCmd = &cobra.Command{
		Use : "lookup",
		Aliases : []string{"lu"},
		Short: "lookup hostnames, IP addresses, MX records, and nameservers.",
		Run : func(cmd *cobra.Command, args []string) {
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
	lookUpCmd.Flags().StringVarP(&ipAddress, "addresses", "a", "", "look up IP addresses by hostname")
	lookUpCmd.Flags().StringVarP(&nameServer, "nameservers", "n", "", "look up name server by hostname")
	lookUpCmd.Flags().StringVar(&hostName, "hostnames", "", "look up hostnames by IP address")
	lookUpCmd.Flags().StringVarP(&domain, "mx", "m", "", "look up MX records by domain")
	Networker.AddCommand(lookUpCmd)
}