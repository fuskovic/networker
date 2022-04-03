package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"cdr.dev/coder-cli/pkg/tablewriter"
	"github.com/fuskovic/networker/internal/list"
	"github.com/fuskovic/networker/internal/ports"
	"github.com/fuskovic/networker/internal/resolve"
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
		if len(os.Args[1:]) == 0 {
			devices, err := list.Devices(ctx)
			if err != nil {
				usage.Fatalf(cmd, "failed to list network devices: %s", err)
			}
			for i := range devices {
				hosts = append(hosts, devices[i].LocalIP.String())
			}
		} else {
			host := os.Args[2]
			ip := net.ParseIP(host)
			if ip == nil {
				if _, err := resolve.AddrByHostName(host); err != nil {
					usage.Fatalf(cmd, "failed to resolve ip address from hostname: %s", err)
				}
			}
			hosts = append(hosts, os.Args[2])
		}

		start := time.Now()
		ticker := time.NewTicker(500 * time.Millisecond)
		done := make(chan bool)

		if shouldOutputAsJSON {
			go func() {
				var dots string
				for {
					select {
					case <-done:
						fmt.Print("\r\n")
					case <-ticker.C:
						dots += "."
						fmt.Printf("\r")
						fmt.Printf("scanning%s", dots)
					}
				}
			}()
		}

		scans, err := ports.NewScanner(hosts, shouldScanAll).Scan(ctx)
		if err != nil {
			usage.Fatalf(cmd, "failed scan hosts: %s", err)
		}

		if !shouldOutputAsJSON {
			ticker.Stop()
			done <- true
			fmt.Printf("\nscan completed in %s\n", time.Since(start).Round(time.Second))
		} else {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "\t")
			enc.SetEscapeHTML(false)
			if err := enc.Encode(scans); err != nil {
				usage.Fatalf(cmd, "failed to encode scan as json: %s", err)
			}
			return
		}

		if err := tablewriter.WriteTable(os.Stdout, len(scans), func(i int) interface{} { return scans[i] }); err != nil {
			usage.Fatalf(cmd, "failed to write scans table: %s", err)
		}
	},
}
