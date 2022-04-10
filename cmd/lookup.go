package cmd

import (
	"net"
	"os"

	"github.com/fuskovic/networker/internal/encoder"
	"github.com/fuskovic/networker/internal/resolve"
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
# Lookup hostname by IP:

	networker lookup hostname 8.8.8.8

# Lookup hostname by IP(short-hand):

	nw lu hn 8.8.8.8

# Lookup hostname by IP(short-hand) and output as json:

	nw lu hn 8.8.8.8 -o json

# Lookup hostname by IP(short-hand) and output as yaml:

	nw lu hn 8.8.8.8 -o yaml

# Lookup IP by hostname:

	networker lookup ip dns.google.

# Lookup IP by hostname(short-hand):

	nw lu ip dns.google.

# Lookup IP by hostname(short-hand) and output as json:

	nw lu ip dns.google. -o json

# Lookup IP by hostname(short-hand) and output as yaml:

	nw lu ip dns.google. -o yaml

# Lookup nameservers by hostname:

	networker lookup nameservers dns.google.

# Lookup nameservers by hostname(short-hand):

	nw lu ns dns.google.

# Lookup nameservers by hostname(short-hand) and output as json:

	nw lu ns dns.google. -o json

# Lookup nameservers by hostname(short-hand) and output as yaml:

	nw lu ns dns.google. -o yaml

# Lookup ISP by ip or hostname:

	networker lookup isp 8.8.8.8
	networker lookup isp dns.google.

# Lookup ISP by ip or hostname(short-hand):

	nw lu isp 8.8.8.8
	nw lu isp dns.google.

# Lookup ISP by ip or hostname(short-hand) and ouput as json:

	nw lu isp 8.8.8.8 -o json
	nw lu isp dns.google. -o json

# Lookup ISP by ip or hostname(short-hand) and ouput as yaml:

	nw lu isp 8.8.8.8 -o yaml
	nw lu isp dns.google. -o yaml

# Lookup network by ip or hostname:

	networker lookup network 8.8.8.8
	networker lookup network dns.google.

# Lookup network by ip or hostname(short-hand):

	nw lu n 8.8.8.8
	nw lu n dns.google.

# Lookup network by ip or hostname(short-hand) and output as json:

	nw lu n 8.8.8.8 -o json
	nw lu n dns.google. -o json

# Lookup network by ip or hostname(short-hand) and output as json:

	nw lu n 8.8.8.8 -o yaml
	nw lu n dns.google. -o yaml

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
	Example: `
# Lookup hostname by IP:

	networker lookup hostname 8.8.8.8

# Lookup hostname by IP(short-hand):

	nw lu hn 8.8.8.8

# Lookup hostname by IP(short-hand) and output as json:

	nw lu hn 8.8.8.8 -o json

# Lookup hostname by IP(short-hand) and output as yaml:

	nw lu hn 8.8.8.8 -o yaml

	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ipAddr := net.ParseIP(args[0])
		if ipAddr == nil {
			usage.Fatalf(cmd, "%q is not a valid ip address", args[0])
		}

		record, err := resolve.HostNameByIP(ipAddr)
		if err != nil {
			usage.Fatalf(cmd, "lookup failed: %s", err)
		}

		enc := encoder.New[resolve.Record](os.Stdout, output)
		if err := enc.Encode(*record); err != nil {
			usage.Fatalf(cmd, "failed to encode hostname record: %s", err)
		}
	},
}

var lookupIpaddressCmd = &cobra.Command{
	Use:   "ip",
	Short: "Lookup the ip address of the provided hostname.",
	Example: `
# Lookup IP by hostname:

	networker lookup ip dns.google.

# Lookup IP by hostname(short-hand):

	nw lu ip dns.google.

# Lookup IP by hostname(short-hand) and output as json:

	nw lu ip dns.google. -o json

# Lookup IP by hostname(short-hand) and output as yaml:

	nw lu ip dns.google. -o yaml

	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if ip := net.ParseIP(args[0]); ip != nil {
			usage.Fatal(cmd, "expected a hostname not an ip address")
			return
		}

		record, err := resolve.AddrByHostName(args[0])
		if err != nil {
			usage.Fatalf(cmd, "lookup failed: %s", err)
		}

		enc := encoder.New[resolve.Record](os.Stdout, output)
		if err := enc.Encode(*record); err != nil {
			usage.Fatalf(cmd, "failed to encode ip address record: %s", err)
		}
	},
}

var lookupIspCmd = &cobra.Command{
	Use: "isp",
	Example: `
# Lookup ISP by ip or hostname:

	networker lookup isp 8.8.8.8
	networker lookup isp dns.google.

# Lookup ISP by ip or hostname(short-hand):

	nw lu isp 8.8.8.8
	nw lu isp dns.google.

# Lookup ISP by ip or hostname(short-hand) and ouput as json:

	nw lu isp 8.8.8.8 -o json
	nw lu isp dns.google. -o json

# Lookup ISP by ip or hostname(short-hand) and ouput as yaml:

	nw lu isp 8.8.8.8 -o yaml
	nw lu isp dns.google. -o yaml
`,
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

		enc := encoder.New[resolve.InternetServiceProvider](os.Stdout, output)
		if err := enc.Encode(*isp); err != nil {
			usage.Fatalf(cmd, "failed to encode internet service provider: %s", err)
		}
	},
}

var lookupNameserversCmd = &cobra.Command{
	Use:     "nameservers",
	Aliases: []string{"ns"},
	Short:   "Lookup nameservers for the provided hostname.",
	Example: `
# Lookup nameservers by hostname:

	networker lookup nameservers dns.google.

# Lookup nameservers by hostname(short-hand):

	nw lu ns dns.google.

# Lookup nameservers by hostname(short-hand) and output as json:

	nw lu ns dns.google. -o json

# Lookup nameservers by hostname(short-hand) and output as yaml:

	nw lu ns dns.google. -o yaml
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hostname, _, err := resolve.HostAndAddr(args[0])
		if err != nil {
			usage.Fatalf(cmd, "%q is an invalid host: %s", args[0], err)
		}

		nameservers, err := resolve.NameServersByHostName(hostname)
		if err != nil {
			usage.Fatalf(cmd, "lookup failed: %s", err)
		}

		enc := encoder.New[resolve.NameServer](os.Stdout, output)
		if err := enc.Encode(nameservers...); err != nil {
			usage.Fatalf(cmd, "failed to encode nameservers: %s", err)
		}
	},
}

var lookupNetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "Lookup the network address of a provided host.",
	Example: `
# Lookup network by ip or hostname:

	networker lookup network 8.8.8.8
	networker lookup network dns.google.

# Lookup network by ip or hostname(short-hand):

	nw lu n 8.8.8.8
	nw lu n dns.google.

# Lookup network by ip or hostname(short-hand) and output as json:

	nw lu n 8.8.8.8 -o json
	nw lu n dns.google. -o json

# Lookup network by ip or hostname(short-hand) and output as json:

	nw lu n 8.8.8.8 -o yaml
	nw lu n dns.google. -o yaml

`,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"n"},
	Run: func(cmd *cobra.Command, args []string) {
		nwRecord, err := resolve.NetworkByHost(args[0])
		if err != nil {
			usage.Fatalf(cmd, "lookup failed: %s", err)
		}

		enc := encoder.New[resolve.NetworkRecord](os.Stdout, output)
		if err := enc.Encode(*nwRecord); err != nil {
			usage.Fatalf(cmd, "failed to encode network record: %s", err)
		}
	},
}
