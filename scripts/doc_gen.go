package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/fuskovic/networker/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var projectRoot string

func init() {
	output, _ := exec.Command("git", "rev-parse", "--show-toplevel").CombinedOutput()
	projectRoot = strings.Replace(string(output), "\n", "", 1)
}

var cmdToDocPathMap = map[*cobra.Command]string{
	cmd.Root:                 path.Join(projectRoot, "docs/networker.md"),
	cmd.ListCmd:              path.Join(projectRoot, "docs/networker_list.md"),
	cmd.LookupCmd:            path.Join(projectRoot, "docs/networker_lookup.md"),
	cmd.LookupHostnameCmd:    path.Join(projectRoot, "docs/networker_lookup_hostname.md"),
	cmd.LookupIpaddressCmd:   path.Join(projectRoot, "docs/networker_lookup_ip.md"),
	cmd.LookupIspCmd:         path.Join(projectRoot, "docs/networker_lookup_isp.md"),
	cmd.LookupNetworkCmd:     path.Join(projectRoot, "docs/networker_lookup_network.md"),
	cmd.LookupNameserversCmd: path.Join(projectRoot, "docs/networker_lookup_nameservers.md"),
	cmd.RequestCmd:           path.Join(projectRoot, "docs/networker_request.md"),
	cmd.ScanCmd:              path.Join(projectRoot, "docs/networker_scan.md"),
}

func main() {
	for cmd, docPath := range cmdToDocPathMap {
		f, err := os.Create(docPath)
		if err != nil {
			log.Fatalf("create %s: %s\n", docPath, err)
		}
		defer f.Close()

		if err := doc.GenMarkdown(cmd, f); err != nil {
			log.Fatalf("generate ")
		}
	}
	log.Println("docs successfully updated")
}
