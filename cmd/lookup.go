package cmd

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/spf13/cobra"
)

type lookUpFunc func(string)

var (
	hostName, ipAddress, nameServer, domain string
	
	supportedLookUps = map[string]lookUpFunc{
		hostName : lookUpHostNamesByIP,
		ipAddress: lookUpIPsByHostName,
		nameServer:lookUpNameServerByHostName,
		domain:	   lookUpMxRecordsForDomain,
	}

	lookup = &cobra.Command{
		Use : "lookup",
		Aliases : []string{"lu"},
		Long: "lookup hostnames, IP addresses, MX records, and nameservers.",
		Run : func(cmd *cobra.Command, args []string) {
			for value, lookUp := range supportedLookUps {
				if value != "" {
					lookUp(value)
				}
			}
		},
	}
)

func init() {
	list.Flags().StringVarP(&ipAddress, "addresses", "a", "", "look up IP addresses by hostname")
	list.Flags().StringVarP(&nameServer, "nameservers", "n", "", "look up name server by hostname")
	list.Flags().StringVarP(&hostName, "hostnames", "h", "", "look up hostnames by IP address")
	list.Flags().StringVarP(&domain, "mx", "m", "", "look up MX records by domain")
	Networker.AddCommand(lookup)
}

func lookUpHostNamesByIP(IP string) {
	if net.ParseIP(IP) == nil {
		log.Fatalf("Valid IP not detected. Value provided: %s\n", IP)
	}

	fmt.Printf("Looking up hostnames for IP address: %s\n", IP)
	hostnames, err := net.LookupAddr(IP)
	if err != nil {
		log.Fatalf("failed to lookup hostnames for %s\nerror : %v\n", IP, err)
	}
	for _, hostnames := range hostnames {
		fmt.Println(hostnames)
	}
}

func lookUpIPsByHostName(hostName string) {
	fmt.Printf("Looking up IP addresses for hostname: %s\n", hostName)
	IPs, err := net.LookupHost(hostName)
	if err != nil {
		log.Fatalf("failed to look up IP addresses for %s\nerror : %v\n", hostName, err)
	}
	for _, ip := range IPs {
		fmt.Println(ip)
	}
}

func lookUpNameServerByHostName(hostName string) {
	fmt.Printf("Looking up nameservers for %s\n", hostName)
	nameservers, err := net.LookupNS(hostName)
	if err != nil {
		log.Fatalf("failed to look up name server for %s\nerror : %v\n", hostName, err)
	}
	for _, ns := range nameservers {
		fmt.Println(ns.Host)
	}
}

func lookUpMxRecordsForDomain(domain string){
	fmt.Printf("Looking up MX records for %s\n", domain)
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Fatalf("Failed to lookup mx records for %s\nerror : %v\n", domain, err)
	}

	for _, mxRecord := range mxRecords {
		fmt.Printf("Host: %s\tPreference: %d\n", mxRecord.Host, mxRecord.Pref)
	}
}

func trim(s string) string {
	return strings.TrimSpace(s)
}