package usage

import (
	"log"

	"github.com/spf13/pflag"
)

// Fatal prints the usage for the flagset and the args before returning an exit code 1.
func Fatal(fl *pflag.FlagSet, args ...interface{}) {
	fl.Usage()
	log.Fatal(args...)
}

// Fatalf prints the usage for the flagset and the formatted args before returning an exit code 1.
func Fatalf(fl *pflag.FlagSet, format string, args ...interface{}) {
	fl.Usage()
	log.Fatalf(format, args...)
}
