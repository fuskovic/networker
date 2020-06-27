package main

import (
	"github.com/fuskovic/networker/cmd"
	"go.coder.com/cli"
)

func main() {
	cli.RunRoot(cmd.RootCmd{})
}
