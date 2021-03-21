package main

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
	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"golang.org/x/xerrors"
)

type scanCmd struct {
	host string
	all  bool
	json bool
}

func (cmd *scanCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "scan",
		Usage:   "[flags]",
		Aliases: []string{"s"},
		Desc:    "Scan hosts for open ports.",
	}
}

func (cmd *scanCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVar(&cmd.host, "host", "", "Host to scan(scans all hosts on LAN if not provided).")
	fl.BoolVarP(&cmd.all, "all", "a", false, "Scan all ports(scans first 1024 if not enabled).")
	fl.BoolVar(&cmd.json, "json", false, "Output as json.")
}

func (cmd *scanCmd) Run(fl *pflag.FlagSet) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hosts, err := cmd.getHostsToScan(ctx)
	if err != nil {
		fl.Usage()
		log.Fatalf("failed to get hosts to scan: %s", err)
	}

	start := time.Now()
	log.Println("scanning...")

	scans, err := ports.NewScanner(hosts, cmd.all).Scan(ctx)
	if err != nil {
		fl.Usage()
		log.Fatalf("failed scan hosts: %s", err)
	}

	log.Printf("scan completed in %s", time.Since(start))

	if cmd.json {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "\t")
		enc.SetEscapeHTML(false)
		if err := enc.Encode(scans); err != nil {
			fl.Usage()
			log.Fatalf("failed to encode scan as json: %s", err)
		}
		return
	}

	if err := tablewriter.WriteTable(os.Stdout, len(scans), func(i int) interface{} { return scans[i] }); err != nil {
		fl.Usage()
		log.Fatalf("failed to write scans table: %s", err)
	}
}

// if the user did not provide a host, get all the hosts on the network.
func (cmd *scanCmd) getHostsToScan(ctx context.Context) ([]string, error) {
	if cmd.host == "" {
		var hosts []string
		devices, err := list.Devices(ctx)
		if err != nil {
			return nil, xerrors.Errorf("failed to list devices: %w", err)
		}
		for i := range devices {
			hosts = append(hosts, devices[i].LocalIP.String())
		}
		return hosts, nil
	}

	if !isValidHost(cmd.host) {
		return nil, xerrors.Errorf("%q is not a valid host", cmd.host)
	}
	return []string{cmd.host}, nil
}

func isValidHost(host string) bool {
	// otherwise, the user may have either provided a hostname or an ip address.
	ip := net.ParseIP(host)
	//if the ip address parse failed, we can assume the user provided a hostname.
	if ip == nil {
		// if we can't look up an ip address by the provided hostname, we can assume the hostname provided was invalid.
		if _, err := resolve.AddrByHostName(host); err != nil {
			return false
		}
	}
	return true
}
