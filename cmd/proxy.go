package cmd

import (
	"fmt"

	"github.com/fuskovic/networker/pkg/proxy"
	"github.com/spf13/cobra"
)

var (
	proxyCfg = &proxy.Config{}

	proxyCmd = &cobra.Command{
		Use:     "proxy",
		Aliases: []string{"p"},
		Example: proxyExample,
		Short:   "Forward network traffic from one network connection to another.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := proxy.Run(proxyCfg); err != nil {
				fmt.Println(err)
				cmd.Usage()
			}
		},
	}
)

func init() {
	proxyCmd.Flags().StringVarP(&proxyCfg.UpStream, "upstream", "u", proxyCfg.UpStream, "Address of server to forward traffic to.")
	proxyCmd.Flags().IntVarP(&proxyCfg.ListenOn, "listen-on", "l", proxyCfg.ListenOn, "Port to listen on.")
	proxyCmd.MarkFlagRequired("upstream")
	proxyCmd.MarkFlagRequired("listen-on")
	Networker.AddCommand(proxyCmd)
}
