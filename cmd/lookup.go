package cmd

import (
	"encoding/json"
	"net"
	"os"

	"github.com/fuskovic/networker/internal/resolve"
	"github.com/fuskovic/networker/internal/table"
	"github.com/fuskovic/networker/internal/usage"
	"github.com/spf13/cobra"
)

func init() {
	lookupCmd.AddCommand(lookupIspCmd)
	lookupCmd.AddCommand(lookupNameserversCmd)
	lookupCmd.AddCommand(lookupHostnameCmd)
	lookupCmd.AddCommand(lookupIpaddressCmd)
	lookupCmd.AddCommand(lookupNetworkCmd)
	Root.AddCommand(lookupCmd)
}

var lookupCmd = &cobra.Command{
	Use:        "lookup",
	Aliases:    []string{"lu"},
	SuggestFor: []string{},
	Example: `
	Lookup IP by hostname:
		networker lookup ip dns.google.

	Lookup hostname by IP:
		networker lookup hostname 8.8.8.8

	Lookup nameservers by hostname:
		networker lookup nameservers dns.google.

	Lookup ISP by ip or hostname:
		networker lookup isp 8.8.8.8
		networker lookup isp dns.google.

	Lookup network by ip or hostname:
		networker lookup network 8.8.8.8
		networker lookup network dns.google.
`,
	Short: "Lookup hostnames, IPs, ISPs, nameservers, and networks.",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Usage()
	},
}

var lookupHostnameCmd = &cobra.Command{
	Use:     "hostname",
	Short:   "Lookup the hostname for a provided ip address.",
	Aliases: []string{"hn"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ipAddr := net.ParseIP(args[0])
		if ipAddr == nil {
			usage.Fatalf(cmd, "%q is not a valid ip address", args[0])
		}

		record, err := resolve.HostNameByIP(ipAddr)
		if err != nil {
			usage.Fatalf(cmd, "lookup failed: %s", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "\t")
			if err := enc.Encode(record); err != nil {
				usage.Fatalf(cmd, "failed to encode hostname record as json: %s", err)
			}
			return
		}

		tw := table.NewWriter(os.Stdout, []resolve.Record{*record})
		if _, err := tw.Write(nil); err != nil {
			usage.Fatalf(cmd, "failed to write hostname table: %s", err)
		}
	},
}

var lookupIpaddressCmd = &cobra.Command{
	Use:   "ip",
	Short: "Lookup the ip address of the provided hostname.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if ip := net.ParseIP(args[0]); ip != nil {
			usage.Fatal(cmd, "expected a hostname not an ip address")
			return
		}

		record, err := resolve.AddrByHostName(args[0])
		if err != nil {
			usage.Fatalf(cmd, "lookup failed: %s", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "\t")
			if err := enc.Encode(record); err != nil {
				usage.Fatalf(cmd, "failed to encode ip address record as json: %s", err)
			}
			return
		}

		tw := table.NewWriter(os.Stdout, []resolve.Record{*record})
		if _, err := tw.Write(nil); err != nil {
			usage.Fatalf(cmd, "failed to write table for hostname record: %s", err)
		}
	},
}

var lookupIspCmd = &cobra.Command{
	Use:   "isp",
	Short: "Lookup the internet service provider of a remote host.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, ip, err := resolve.HostAndAddr(args[0])
		if err != nil {
			usage.Fatalf(cmd, "%q is an invalid host: %s", args[0], err)
		}

		if resolve.IsPrivate(ip) {
			usage.Fatalf(cmd, "cannot retrieve internet service provider for private ip")
		}

		isp, err := resolve.ServiceProvider(ip)
		if err != nil {
			usage.Fatalf(cmd, "failed to resolve internet service provider for %q: %s", ip, err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "\t")
			if err := enc.Encode(isp); err != nil {
				usage.Fatalf(cmd, "failed to encode internet service provider as json: %s", err)
			}
			return
		}

		tw := table.NewWriter(os.Stdout, []resolve.InternetServiceProvider{*isp})
		if _, err := tw.Write(nil); err != nil {
			usage.Fatalf(cmd, "failed to write table for service provider record: %s", err)
		}
	},
}

var lookupNameserversCmd = &cobra.Command{
	Use:     "nameservers",
	Aliases: []string{"ns"},
	Short:   "Lookup nameservers for the provided hostname.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hostname, _, err := resolve.HostAndAddr(args[0])
		if err != nil {
			usage.Fatalf(cmd, "%q is an invalid host: %s", args[0], err)
		}

		nameservers, err := resolve.NameServersByHostName(hostname)
		if err != nil {
			usage.Fatalf(cmd, "lookup failed: %s", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "\t")
			if err := enc.Encode(nameservers); err != nil {
				usage.Fatalf(cmd, "failed to encode nameservers as json: %s", err)
			}
			return
		}

		tw := table.NewWriter(os.Stdout, nameservers)
		if _, err := tw.Write(nil); err != nil {
			usage.Fatalf(cmd, "failed to write table for nameservers: %s", err)
		}
	},
}

var lookupNetworkCmd = &cobra.Command{
	Use:     "network",
	Short:   "Lookup the network address of a provided host.",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"n"},
	Run: func(cmd *cobra.Command, args []string) {
		record, err := resolve.NetworkByHost(args[0])
		if err != nil {
			usage.Fatalf(cmd, "lookup failed: %s", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "\t")
			if err := enc.Encode(record); err != nil {
				usage.Fatalf(cmd, "failed to encode network record as json: %s", err)
			}
			return
		}

		tw := table.NewWriter(os.Stdout, []resolve.Record{*record})
		if _, err := tw.Write(nil); err != nil {
			usage.Fatalf(cmd, "failed to write table for network record: %s", err)
		}
	},
}
