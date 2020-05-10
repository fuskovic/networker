package cmd

import (
	"fmt"
	"log"

	"github.com/fuskovic/networker/pkg/lookup"
	"github.com/spf13/cobra"
)

type lookUpFunc func(string) error

var (
	hostName, ipAddress, nameServer, domain, network string
	networkEx                                        = "networker lookup --network www.farishuskovic.dev"
	hostNameEx                                       = "networker lookup --hostnames 192.81.212.192"
	nameServerEx                                     = "networker lookup --nameservers farishuskovic.dev"
	ipEx                                             = "networker lookup --addresses farishuskovic.dev"
	lookUpExFormat                                   = "\nlookup network : %s\nlookup hostname : %s\nlookup nameserver : %s\nlookup ip : %s\n"
	lookUpEx                                         = fmt.Sprintf(lookUpExFormat, networkEx, hostNameEx, nameServerEx, ipEx)

	supportedLookUps = map[*string]lookUpFunc{
		&hostName:   lookup.HostNamesByIP,
		&ipAddress:  lookup.IPsByHostName,
		&nameServer: lookup.NameServersByHostName,
		&domain:     lookup.MxRecordsForDomain,
		&network:    lookup.NetworkByHostName,
	}

	lookUpCmd = &cobra.Command{
		Use:     "lookup",
		Aliases: []string{"lu"},
		Example: lookUpEx,
		Short:   "lookup hostnames, IP addresses, MX records, and nameservers.",
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
	lookUpCmd.Flags().StringVar(&network, "network", "", "look up the network for a hostname")
	lookUpCmd.Flags().StringVarP(&ipAddress, "addresses", "a", "", "look up IP addresses by hostname")
	lookUpCmd.Flags().StringVarP(&nameServer, "nameservers", "n", "", "look up name server by hostname")
	lookUpCmd.Flags().StringVar(&hostName, "hostnames", "", "look up hostnames by IP address")
	lookUpCmd.Flags().StringVarP(&domain, "mx", "m", "", "look up MX records by domain")
	Networker.AddCommand(lookUpCmd)
}
