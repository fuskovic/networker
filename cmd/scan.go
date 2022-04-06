package cmd

import (
	"context"
	"encoding/json"
	"net"
	"os"

	"github.com/fuskovic/networker/internal/list"
	"github.com/fuskovic/networker/internal/ports"
	"github.com/fuskovic/networker/internal/resolve"
	"github.com/fuskovic/networker/internal/table"
	"github.com/fuskovic/networker/internal/usage"
	"github.com/spf13/cobra"
)

var shouldScanAll, shouldOutputAsJSON bool

func init() {
	scanCmd.Flags().BoolVar(&shouldScanAll, "all-ports", false, "Scan all ports(scans first 1024 if not enabled).")
	scanCmd.Flags().BoolVar(&shouldOutputAsJSON, "json", false, "Output as json.")
	Root.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{
	Use:     "scan",
	Aliases: []string{"s"},
	Short:   "Scan hosts for open ports.",
	Example: `
	Scan well-known ports of single device on network:
		networker scan 127.0.0.1

	Scan well-known ports of all devices on network:
		networker scan

	Scan all ports of single device on network:
		networker scan 127.0.0.1 --all-ports

	Output a scan as json:
		networker scan 127.0.0.1 --json
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
				addr, err := resolve.AddrByHostName(args[0])
				if err != nil {
					usage.Fatalf(cmd, "failed to resolve ip address from hostname: %s", err)
				}
				ip = *addr
			}
			hosts = append(hosts, ip.String())
		}

		scans, err := ports.NewScanner(hosts, shouldScanAll).Scan(ctx)
		if err != nil {
			usage.Fatalf(cmd, "failed scan hosts: %s", err)
		}

		if shouldOutputAsJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "\t")
			if err := enc.Encode(scans); err != nil {
				usage.Fatalf(cmd, "failed to encode scan as json: %s", err)
			}
			return
		}

		tableWriter := table.NewWriter(os.Stdout, scans)
		if _, err := tableWriter.Write(nil); err != nil {
			usage.Fatalf(cmd, "failed to write scans table: %s", err)
		}
	},
}
