package cmd

import (
	"fmt"
	"log"

	"github.com/fuskovic/networker/pkg/lookup"
	"github.com/spf13/cobra"
)

type lookUpFunc func(string) error

var (
	hostName, ipAddress, nameServer, network string
	networkEx                                = "networker lookup --network facebook.com || 31.13.65.36\n"
	hostNameEx                               = "networker lookup --hostnames 157.240.195.35\n"
	nameServerEx                             = "networker lookup --nameservers youtube.com\n"
	ipEx                                     = "networker lookup --addresses youtube.com\n"
	lookUpExFormat                           = "\nlookup network : \n%s\nlookup hostname : \n%s\nlookup nameserver : \n%s\nlookup ip : \n%s\n"
	lookUpEx                                 = fmt.Sprintf(lookUpExFormat, networkEx, hostNameEx, nameServerEx, ipEx)

	supportedLookUps = map[*string]lookUpFunc{
		&hostName:   lookup.HostNamesByIP,
		&ipAddress:  lookup.IPsByHostName,
		&nameServer: lookup.NameServersByHostName,
		&network:    lookup.NetworkByHostName,
	}

	lookUpCmd = &cobra.Command{
		Use:     "lookup",
		Aliases: []string{"lu"},
		Example: lookUpEx,
		Short:   "lookup hostnames, IP addresses, nameservers, and general network information.",
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
	lookUpCmd.Flags().StringVar(&network, "network", "", "look up the network of a host")
	lookUpCmd.Flags().StringVarP(&ipAddress, "addresses", "a", "", "look up IP addresses by hostname")
	lookUpCmd.Flags().StringVarP(&nameServer, "nameservers", "n", "", "look up name servers by hostname")
	lookUpCmd.Flags().StringVar(&hostName, "hostnames", "", "look up hostnames by IP address")
	Networker.AddCommand(lookUpCmd)
}
