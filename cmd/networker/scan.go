package networker

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"
	"time"

	"cdr.dev/coder-cli/pkg/tablewriter"
	"github.com/fuskovic/networker/internal/list"
	"github.com/fuskovic/networker/internal/ports"
	"github.com/fuskovic/networker/internal/resolve"
	"github.com/fuskovic/networker/internal/usage"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
)

type scanCmd struct {
	shouldScanAllPorts bool
	json               bool
}

func (cmd *scanCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "scan",
		Usage:   "[flags] [host]",
		Aliases: []string{"s"},
		Desc:    "Scan hosts for open ports.",
	}
}

func (cmd *scanCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.BoolVarP(&cmd.shouldScanAllPorts, "all", "a", false, "Scan all ports(scans first 1024 if not enabled).")
	fl.BoolVar(&cmd.json, "json", false, "Output as json.")
}

func (cmd *scanCmd) Run(fl *pflag.FlagSet) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var hosts []string
	if len(os.Args) < 3 {
		devices, err := list.Devices(ctx)
		if err != nil {
			usage.Fatalf(fl, "failed to list network devices: %s", err)
		}
		for i := range devices {
			hosts = append(hosts, devices[i].LocalIP.String())
		}
	} else {
		host := os.Args[2]
		ip := net.ParseIP(host)
		if ip == nil {
			if _, err := resolve.AddrByHostName(host); err != nil {
				usage.Fatalf(fl, "failed to resolve ip address from hostname: %s", err)
			}
		}
		hosts = append(hosts, os.Args[2])
	}

	start := time.Now()
	log.Println("scanning...")

	scans, err := ports.NewScanner(hosts, cmd.shouldScanAllPorts).Scan(ctx)
	if err != nil {
		usage.Fatalf(fl, "failed scan hosts: %s", err)
	}

	log.Printf("scan completed in %s", time.Since(start))

	if cmd.json {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "\t")
		enc.SetEscapeHTML(false)
		if err := enc.Encode(scans); err != nil {
			usage.Fatalf(fl, "failed to encode scan as json: %s", err)
		}
		return
	}

	if err := tablewriter.WriteTable(os.Stdout, len(scans), func(i int) interface{} { return scans[i] }); err != nil {
		usage.Fatalf(fl, "failed to write scans table: %s", err)
	}
}
