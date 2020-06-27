package cmd

import (
	"fmt"
	"net"
	"strings"

	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"
)

type (
	lookUpCmd struct {
		hostName, ipAddress, nameServer, network string
	}

	lookUpFunc func(string) error
)

// Spec returns a command spec containing a description of it's usage.
func (cmd *lookUpCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "lookup",
		Usage: "TODO: ADD USAGE",
		Desc:  "Lookup hostnames, IP addresses, and nameservers",
	}
}

// RegisterFlags initializes how a flag set is processed for a particular command.
func (cmd *lookUpCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVarP(&cmd.network, "network", "n", "", "Look up the network a given hostname belongs to.")
	fl.StringVarP(&cmd.ipAddress, "addresses", "a", "", "Look up IP addresses for a given hostname.")
	fl.StringVarP(&cmd.nameServer, "nameservers", "s", "", "Look up nameservers for a given hostname.")
	fl.StringVar(&cmd.hostName, "hostnames", "", "Look up hostnames for a given IP address.")
}

// Run iterates over the flagset in search of supported lookups and prints the information requested.
func (cmd *lookUpCmd) Run(fl *pflag.FlagSet) {
	for value, lookUp := range cmd.supportedLookUps() {
		if value != "" {
			if err := lookUp(value); err != nil {
				flog.Fatal(err.Error())
			}
		}
	}
}

func (cmd *lookUpCmd) supportedLookUps() map[string]lookUpFunc {
	return map[string]lookUpFunc{
		cmd.hostName:   hostNamesByIP,
		cmd.ipAddress:  addrsByHostName,
		cmd.nameServer: nameServersByHostName,
		cmd.network:    networkByHostName,
	}
}

func hostNamesByIP(addr string) error {
	addr = trim(addr)
	flog.Info("Looking up hostnames for IP address: %s", addr)

	if net.ParseIP(addr) == nil {
		return fmt.Errorf("%s is not a valid IP address", addr)
	}

	hostnames, err := net.LookupAddr(addr)
	if err != nil {
		return fmt.Errorf("failed to lookup hostnames for %s\nerror : %v", addr, err)
	}

	if len(hostnames) == 0 {
		return fmt.Errorf("no hostnames found")
	}

	for _, hn := range hostnames {
		flog.Info(hn)
	}
	return nil
}

func addrsByHostName(hostName string) error {
	hostName = trim(hostName)
	flog.Info("Looking up IP addresses for hostname: %s", hostName)

	addrs, err := net.LookupHost(hostName)
	if err != nil {
		return fmt.Errorf("failed to look up IP addresses for %s\nerror : %v", hostName, err)
	}

	if len(addrs) == 0 {
		return fmt.Errorf("no IP addresses found")
	}

	for _, a := range addrs {
		flog.Info(a)
	}
	return nil
}

func nameServersByHostName(hostName string) error {
	hostName = trim(hostName)
	flog.Info("Looking up nameservers for %s", hostName)

	nameservers, err := net.LookupNS(hostName)
	if err != nil {
		return fmt.Errorf("failed to look up name server for %s\nerror : %v", hostName, err)
	}

	if len(nameservers) == 0 {
		return fmt.Errorf("no name servers found")
	}

	for _, ns := range nameservers {
		flog.Info(ns.Host)
	}
	return nil
}

func networkByHostName(hostName string) error {
	ip, err := net.ResolveIPAddr("ip", hostName)
	if err != nil {
		return fmt.Errorf("failed to resolve IP address from hostname : %s\nerror : %v", hostName, err)
	}

	addr := net.ParseIP(ip.String())
	if addr == nil {
		return fmt.Errorf("failed to validate the resolved IP : %s for hostname : %s", addr, hostName)
	}

	mask := addr.DefaultMask()
	network := addr.Mask(mask)
	format := "Hostname : %s\nAddress : %s\nNetwork : %s\n"
	flog.Info(format, hostName, addr.String(), network)
	return nil
}

func trim(s string) string {
	return strings.TrimSpace(s)
}
