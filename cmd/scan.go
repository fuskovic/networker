package cmd

import (
	"context"
	"net"
	"os"
	"slices"

	"github.com/spf13/cobra"

	"github.com/fuskovic/networker/v3/internal/encoder"
	"github.com/fuskovic/networker/v3/internal/list"
	"github.com/fuskovic/networker/v3/internal/resolve"
	"github.com/fuskovic/networker/v3/internal/scanner"
	"github.com/fuskovic/networker/v3/internal/spinner"
	"github.com/fuskovic/networker/v3/internal/usage"
)

var scanAllPorts bool

func init() {
	scanCmd.Flags().BoolVar(&scanAllPorts, "all-ports", false, "Scan all ports(scans first 1024 if not enabled).")
	Root.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{
	Use:     "scan",
	Aliases: []string{"s"},
	Short:   "Scan hosts for open ports.",
	Example: `
# Scan well-known ports(first 1024) of all devices on network:

		networker scan

# Scan well-known ports(first 1024) of all devices on network(short-hand):

		nw s

# Scan well-known ports(first 1024) of all devices on network(short-hand) and output as json:

		nw s -o json

# Scan well-known ports(first 1024) of all devices on network(short-hand) and output as yaml:

		nw s -o yaml

# Scan all ports of all devices on network:

		networker scan --all-ports

# Scan all ports of all devices on network(short-hand):

		nw s --all-ports

# Scan all ports of all devices on network(short-hand) and output as json:

		nw s -o json --all-ports

# Scan all ports of all devices on network(short-hand) and output as yaml:

		nw s -o yaml --all-ports

# Scan well-known ports(first 1024) of single host:

		networker scan localhost

# Scan well-known ports(first 1024) of single host(short-hand):

		nw s localhost

# Scan well-known ports(first 1024) of single host(short-hand) and output as json:

		nw s localhost -o json

# Scan well-known ports(first 1024) of single host(short-hand) and output as yaml:

		nw s localhost -o yaml

# Scan all ports of single host:

		networker scan localhost --all-ports

# Scan all ports of single host(short-hand):

		nw s localhost --all-ports

# Scan all ports of single host(short-hand) and output as json:

		nw s localhost -o json --all-ports

# Scan all ports of single host(short-hand) and output as yaml:

		nw s localhost -o yaml --all-ports

`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var hosts []string
		if len(args) == 0 {
			devices, err := list.Devices(ctx)
			if err != nil {
				usage.Fatalf(cmd, "failed to list network devices: %s", err)
			}
			for i := range devices {
				hosts = append(hosts, devices[i].LocalIP.String())
			}
		} else {
			ip := net.ParseIP(args[0])
			if ip == nil {
				record, err := resolve.AddrByHostName(args[0])
				if err != nil {
					usage.Fatalf(cmd, "failed to resolve ip address from hostname: %s", err)
				}
				ip = record.IP
			}
			hosts = append(hosts, ip.String())
		}

		spinner.Start()

		scans, err := scanner.New(hosts, scanAllPorts).Scan(ctx)
		if err != nil {
			usage.Fatalf(cmd, "failed scan hosts: %s", err)
		}

		spinner.Stop()

		scans = slices.DeleteFunc(scans,
			func(s scanner.Scan) bool {
				return s.Host == "N/A" && len(s.Ports) == 0
			},
		)

		enc := encoder.New[scanner.Scan](os.Stdout, output)
		if err := enc.Encode(scans...); err != nil {
			usage.Fatalf(cmd, "failed to encode devices: %s", err)
		}
	},
}
