package main

import (
	"log"
	"os/exec"
	"path"
	"strings"

	"github.com/fuskovic/networker/cmd"
	"github.com/spf13/cobra/doc"
)

var projectRoot string

func init() {
	output, _ := exec.Command("git", "rev-parse", "--show-toplevel").CombinedOutput()
	projectRoot = strings.Replace(string(output), "\n", "", 1)
}

func main() {
	if err := doc.GenMarkdownTree(cmd.Root, path.Join(projectRoot, "docs")); err != nil {
		log.Fatalf("gen markdown tree: %v\n", err)
	}
	log.Println("docs successfully updated")
}
