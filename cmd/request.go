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
		Example: "TODO: add request example",
		Short:   "Send an HTTP GET request or POST JSON.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := request.Run(reqCfg); err != nil {
				fmt.Println(err)
				cmd.Usage()
			}
		},
	}
)

func init() {
	reqCmd.Flags().StringVarP(&reqCfg.URL, "url", "u", reqCfg.URL, "URL to send request.")
	reqCmd.Flags().StringVarP(&reqCfg.Method, "method", "m", reqCfg.Method, "Specify method. (supported methods include GET and POST)")
	reqCmd.Flags().StringVarP(&reqCfg.Data, "data", "d", reqCfg.Data, "JSON string or file to use for request body.")
	reqCmd.MarkFlagRequired("url")
	reqCmd.MarkFlagRequired("method")
	Networker.AddCommand(reqCmd)
}
