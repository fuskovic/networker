package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// Networker is the root cmd.
var Networker = &cobra.Command{
	Use:  "networker",
	Short: "networker is an easy to use networking CLI tool",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Usage(); err != nil {
			log.Fatalf("failed to print usage - %v\n", err)
		}
	},
}

// Execute runs the root cmd.
func Execute() {
	if err := Networker.Execute(); err != nil {
		log.Fatalf("failed to execute networker - err : %v\n", err)
	}
}
