package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fuskovic/networker/pkg/request"
)

var (
	reqCfg = &request.Config{}

	reqCmd = &cobra.Command{
		Use:     "request",
		Aliases: []string{"req", "r"},
		Example: requestExample,
		Short:   "Send an HTTP request.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := request.Run(reqCfg); err != nil {
				fmt.Println(err)
				cmd.Usage()
			}
		},
	}
)

func init() {
	reqCmd.Flags().StringSliceVarP(&reqCfg.Headers, "add-headers", "a", reqCfg.Headers, "Add a list of comma-separated request headers. (format : key:value,key:value,etc...)")
	reqCmd.Flags().StringVarP(&reqCfg.URL, "url", "u", reqCfg.URL, "URL to send request.")
	reqCmd.Flags().StringVarP(&reqCfg.Method, "method", "m", "GET", "Specify method. (supported methods include GET, POST, PUT, PATCH, and DELETE)")
	reqCmd.Flags().StringVarP(&reqCfg.File, "file", "f", reqCfg.File, "Path to JSON or XML file to use for request body. (content-type headers for each file-type are set automatically)")
	reqCmd.Flags().IntVarP(&reqCfg.TimeOut, "time-out", "t", 3, "Specify number of seconds for time-out.")
	reqCmd.MarkFlagRequired("url")
	Networker.AddCommand(reqCmd)
}
