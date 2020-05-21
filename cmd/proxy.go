package cmd

import (
	"github.com/fuskovic/networker/pkg/proxy"
	"github.com/spf13/cobra"
)

var (
	upStream string
	listenOn int

	proxyCmd = &cobra.Command{
		Use:     "proxy",
		Aliases: []string{"p"},
		Example: proxyExample,
		Short:   "Forward network traffic from one network connection to another.",
		Run: func(cmd *cobra.Command, args []string) {
			proxy.Run(listenOn, upStream)
		},
	}
)

func init() {
	proxyCmd.Flags().StringVarP(&upStream, "upstream", "u", upStream, "Address of server to forward traffic to.")
	proxyCmd.Flags().IntVarP(&listenOn, "listen-on", "l", listenOn, "Port to listen on.")
	proxyCmd.MarkFlagRequired("upstream")
	proxyCmd.MarkFlagRequired("listen-on")
	Networker.AddCommand(proxyCmd)
}
