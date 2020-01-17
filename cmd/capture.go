package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/fuskovic/networker/pkg/capture"
	"github.com/google/gopacket/pcapgo"
	"github.com/spf13/cobra"
)

var (
	devices      []string
	seconds      int64
	longExample  = "long form:\n./networker capture --devices en1 --seconds 10 --out myCaptureSession --limit --num 100"
	shortExample = "short form:\n./networker c -d en1 -s 10 -o myCaptureSession -l -n 100"
	outFile      string
	limit        bool
	numToCapture int64
	writer       *pcapgo.Writer
	file         *os.File
	err          error

	captureCmd = &cobra.Command{
		Use:     "capture",
		Aliases: []string{"c", "cap"},
		Short:   "capture network packets on specified devices.",
		Example: fmt.Sprintf("%s\n%s\n", longExample, shortExample),
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

			if err := capture.Packets(devices, seconds, numToCapture, writer, limit); err != nil {
				log.Printf("error during packet capture : %v\n", err)
				cmd.Usage()
				return
			}
		},
	}
)

func init() {
	captureCmd.PersistentFlags().Int64VarP(&seconds, "seconds", "s", 0, "Amount of seconds to run capture")
	captureCmd.PersistentFlags().StringSliceVarP(&devices, "devices", "d", []string{}, "devices on which to capture network packets (comma separated).")
	captureCmd.PersistentFlags().StringVarP(&outFile, "out", "o", "", "specify outfile to write captured packets to")
	captureCmd.PersistentFlags().BoolVarP(&limit, "limit", "l", false, "enable packet capture limiting(must use with --num || -n to specify number).")
	captureCmd.PersistentFlags().Int64VarP(&numToCapture, "num", "n", 0, "number of packets to capture (accumulative for all devices)")
	captureCmd.MarkPersistentFlagRequired("seconds")
	Networker.AddCommand(captureCmd)
}
