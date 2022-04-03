package cmd

import (
	"encoding/json"
	"log"
	"net"
	"os"

	"cdr.dev/coder-cli/pkg/tablewriter"
	"github.com/fuskovic/networker/internal/resolve"
	"github.com/fuskovic/networker/internal/usage"
	"github.com/spf13/cobra"
)

var (
	ip                 string
	hostname           string
	host               string
)

func init() {
	lookupHostnameCmd.Flags().StringVar(&ip, "ip", "", "IP address.")
	_ = lookupHostnameCmd.MarkFlagRequired("ip")
	lookupCmd.AddCommand(lookupHostnameCmd)

	lookupIpaddressCmd.Flags().StringVar(&hostname, "hostname", "", "Hostname.")
	_ = lookupIpaddressCmd.MarkFlagRequired("hostname")
	lookupCmd.AddCommand(lookupIpaddressCmd)

	lookupIspCmd.Flags().StringVar(&host, "host", "", "IP address or hostname to get the network address for.")
	lookupIspCmd.Flags().BoolVar(&shouldOutputAsJSON, "json", false, "Output as JSON.")
	_ = lookupIspCmd.MarkFlagRequired("host")
	lookupCmd.AddCommand(lookupIspCmd)

	lookupNameserversCmd.Flags().StringVar(&hostname, "hostname", "", "Hostname.")
	lookupNameserversCmd.Flags().BoolVar(&shouldOutputAsJSON, "json", false, "Output as JSON.")
	_ = lookupNameserversCmd.MarkFlagRequired("hostname")
	lookupCmd.AddCommand(lookupNameserversCmd)

	lookupNetworkCmd.Flags().StringVar(&host, "host", "", "IP address or hostname to get the network address for.")
	_ = lookupNetworkCmd.MarkFlagRequired("host")
	lookupCmd.AddCommand(lookupNetworkCmd)

	Root.AddCommand(lookupCmd)
}

var lookupCmd = &cobra.Command{
	Use:        "lookup",
	Aliases:    []string{"lu"},
	SuggestFor: []string{},
	Example: `
	Lookup IP by hostname:
		networker lookup ip --hostname dns.google.

	Lookup hostname by IP:
		networker lookup hostname --ip 8.8.8.8

	Lookup nameservers by hostname:
		networker lookup nameservers --hostname dns.google.

	Lookup ISP by ip or hostname:
		networker lookup isp --host 8.8.8.8
		networker lookup network --host dns.google.

	Lookup network by ip or hostname:
		networker lookup network --host 8.8.8.8
		networker lookup network --host dns.google.
`,
	Short: "Lookup hostnames, IPs, ISPs, nameservers, and networks.",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Usage()
	},
}

var lookupHostnameCmd = &cobra.Command{
	Use:   "hostname",
	Short: "Lookup the hostname for a provided ip address.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ipAddr := net.ParseIP(ip)
		if ipAddr == nil {
			usage.Fatalf(cmd, "%q is not a valid ip address", ip)
		}

		hostname, err := resolve.HostNameByIP(ipAddr)
		if err != nil {
			usage.Fatalf(cmd, "lookup failed: %s", err)
		}
		log.Printf("lookup successful - hostname: %s", hostname)
	},
}

var lookupIpaddressCmd = &cobra.Command{
	Use:   "ip",
	Short: "Lookup the ip address of the provided hostname.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ipAddr, err := resolve.AddrByHostName(hostname)
		if err != nil {
			usage.Fatalf(cmd, "lookup failed: %s", err)
		}
		log.Printf("lookup successful - ip-address: %s", ipAddr)
	},
}

var lookupIspCmd = &cobra.Command{
	Use:   "isp",
	Short: "Lookup the internet service provider of a remote host.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		_, ip, err := resolve.HostAndAddr(host)
		if err != nil {
			usage.Fatalf(cmd, "%q is an invalid host: %s", host, err)
		}

		if resolve.IsPrivate(ip) {
			usage.Fatalf(cmd, "cannot retrieve internet service provider for private ip")
		}

		isp, err := resolve.ServiceProvider(ip)
		if err != nil {
			usage.Fatalf(cmd, "failed to resolve internet service provider for %q: %s", host, err)
		}

		if shouldOutputAsJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "\t")
			enc.SetEscapeHTML(false)
			if err := enc.Encode(isp); err != nil {
				usage.Fatalf(cmd, "failed to encode internet service provider as json: %s", err)
			}
			return
		}

		err = tablewriter.WriteTable(os.Stdout, 1,
			func(_ int) interface{} {
				return *isp
			},
		)

		if err != nil {
			usage.Fatalf(cmd, "failed to write service provider table for %q: %s", host, err)
		}
	},
}

var lookupNameserversCmd = &cobra.Command{
	Use:     "nameservers",
	Aliases: []string{"ns"},
	Short:   "Lookup nameservers for the provided hostname.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if ip := net.ParseIP(hostname); ip != nil {
			usage.Fatalf(cmd, "can only lookup nameservers by hostname; not ip.")
		}

		nameservers, err := resolve.NameServersByHostName(hostname)
		if err != nil {
			usage.Fatalf(cmd, "lookup failed: %s", err)
		}

		if shouldOutputAsJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "\t")
			enc.SetEscapeHTML(false)
			if err := enc.Encode(nameservers); err != nil {
				usage.Fatalf(cmd, "failed to encode nameservers as json: %s", err)
			}
			return
		}

		err = tablewriter.WriteTable(os.Stdout, len(nameservers),
			func(i int) interface{} {
				return nameservers[i]
			},
		)

		if err != nil {
			usage.Fatalf(cmd, "failed to write nameservers table: %s", err)
		}
	},
}

var lookupNetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "Lookup the network address of a provided host.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		network, err := resolve.NetworkByHost(host)
		if err != nil {
			usage.Fatalf(cmd, "lookup failed: %s", err)
		}
		log.Printf("lookup successful - network address: %s", network)
	},
}
