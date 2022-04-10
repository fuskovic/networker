package usage

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// Fatal prints the usage for the flagset and the args before returning an exit code 1.
func Fatal(cmd *cobra.Command, args ...interface{}) {
	log.Print(args...)
	_ = cmd.Usage()
	os.Exit(1)
}

// Fatalf prints the usage for the flagset and the formatted args before returning an exit code 1.
func Fatalf(cmd *cobra.Command, format string, args ...interface{}) {
	log.Printf(format, args...)
	_ = cmd.Usage()
	os.Exit(1)
}
