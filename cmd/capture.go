package cmd

import (
	"log"

	"github.com/fuskovic/networker/pkg/capture"
	"github.com/spf13/cobra"
)

var (
	devices []string
	seconds int64
	example = "networker capture eth0 -s 10"

	captureCmd = &cobra.Command{
		Use:     "capture",
		Aliases: []string{"c", "cap"},
		Short:   "capture network packets on specified devices.",
		Example: example,
		Run: func(cmd *cobra.Command, args []string) {
			if seconds < 5 {
				log.Printf("capture must be at least 5 seconds long - your input : %d\n", seconds)
				return
			}

			if err := capture.Packets(devices, seconds); err != nil {
				log.Printf("error during packet capture : %v\n", err)
				return
			}
		},
	}
)

func init() {
	captureCmd.PersistentFlags().Int64VarP(&seconds, "seconds", "s", 0, "Amount of seconds to run capture")
	captureCmd.MarkPersistentFlagRequired("seconds")
	captureCmd.Flags().StringSliceVarP(&devices, "devices", "d", []string{}, "devices on which to capture network packets (comma separated).")
	Networker.AddCommand(captureCmd)
}
