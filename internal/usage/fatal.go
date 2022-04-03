package usage

import (
	"log"

	"github.com/spf13/cobra"
)

// Fatal prints the usage for the flagset and the args before returning an exit code 1.
func Fatal(cmd *cobra.Command, args ...interface{}) {
	cmd.Usage()
	log.Fatal(args...)
}

// Fatalf prints the usage for the flagset and the formatted args before returning an exit code 1.
func Fatalf(cmd *cobra.Command, format string, args ...interface{}) {
	cmd.Usage()
	log.Fatalf(format, args...)
}
