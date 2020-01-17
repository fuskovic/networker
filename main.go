package main

import (
	"log"

	"github.com/fuskovic/networker/cmd"
)

func main() {
	if err := cmd.Networker.Execute(); err != nil {
		log.Fatalf("failed to execute networker\nerror : %v\n", err)
	}
}