package cmd

import (
	"log"
	"os"

	"github.com/fuskovic/networker/pkg/capture"
	"github.com/google/gopacket/pcapgo"
	"github.com/spf13/cobra"
)

var (
	devices        []string
	seconds        int64
	outFile        string
	limit, verbose bool
	numToCapture   int64
	writer         *pcapgo.Writer
	file           *os.File
	err            error

	captureCmd = &cobra.Command{
		Use:     "capture",
		Aliases: []string{"c", "cap"},
		Short:   "Capture network packets on specified devices.",
		Example: captureExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(devices) == 0 {
				log.Printf("no designated devices")
				cmd.Usage()
				return
			}

			if seconds < 5 {
				log.Printf("capture must be at least 5 seconds long - your input : %d\n", seconds)
				cmd.Usage()
				return
			}

			if limit && numToCapture < 1 {
				log.Printf("use of --limit flag without use of --num flag\nPlease specify number of packets to limit capture\nminimum is 1")
				cmd.Usage()
				return
			}

			if outFile != "" {
				file, writer, err = capture.NewWriter(outFile)
				if err != nil {
					log.Printf("failed to create a new writer - err : %v\n", err)
					cmd.Usage()
					return
				}
				defer file.Close()
			}

			if err := capture.Packets(devices, seconds, numToCapture, writer, limit, verbose); err != nil {
				log.Printf("error during packet capture : %v\n", err)
				cmd.Usage()
				return
			}
		},
	}
)

func init() {
	captureCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging.")
	captureCmd.Flags().Int64VarP(&seconds, "seconds", "s", 0, "Amount of seconds to run capture for.")
	captureCmd.Flags().StringSliceVarP(&devices, "devices", "d", []string{}, "Comma-separated list of devices to capture packets on.")
	captureCmd.Flags().StringVarP(&outFile, "out", "o", "", "Name of an output file to write the packets to.")
	captureCmd.Flags().BoolVarP(&limit, "limit", "l", false, "Limit the number of packets to capture. (must be used with the --num flag)")
	captureCmd.Flags().Int64VarP(&numToCapture, "num", "n", 0, "Number of total packets to capture across all devices.")
	captureCmd.MarkFlagRequired("seconds")
	Networker.AddCommand(captureCmd)
}
