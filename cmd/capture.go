package cmd

import (
	"fmt"

	"github.com/fuskovic/networker/pkg/capture"
	"github.com/spf13/cobra"
)

var (
	capCfg = &capture.Config{}

	captureCmd = &cobra.Command{
		Use:     "capture",
		Aliases: []string{"c", "cap"},
		Short:   "Capture network packets on specified devices.",
		Example: captureExample,
		Run: func(cmd *cobra.Command, args []string) {
			if err := capture.Run(capCfg); err != nil {
				fmt.Println(err)
				cmd.Usage()
			}
		},
	}
)

func init() {
	captureCmd.Flags().BoolVarP(&capCfg.Verbose, "verbose", "v", capCfg.Verbose, "Enable verbose logging.")
	captureCmd.Flags().Int64VarP(&capCfg.Seconds, "seconds", "s", capCfg.Seconds, "Amount of seconds to run capture for.")
	captureCmd.Flags().StringSliceVarP(&capCfg.Devices, "devices", "d", capCfg.Devices, "Comma-separated list of devices to capture packets on.")
	captureCmd.Flags().StringVarP(&capCfg.OutFile, "out", "o", capCfg.OutFile, "Name of an output file to write the packets to.")
	captureCmd.Flags().BoolVarP(&capCfg.Limit, "limit", "l", capCfg.Limit, "Limit the number of packets to capture. (must be used with the --num flag)")
	captureCmd.Flags().Int64VarP(&capCfg.NumToCapture, "num", "n", capCfg.NumToCapture, "Number of total packets to capture across all devices.")
	captureCmd.MarkFlagRequired("seconds")
	Networker.AddCommand(captureCmd)
}
