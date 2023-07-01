package main

import "github.com/fuskovic/networker/v3/cmd"

func main() {
	cmd.Root.CompletionOptions.DisableDefaultCmd = true
	cmd.Root.Execute()
}
