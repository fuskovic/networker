package cmd

import (
	"log"

	"github.com/fuskovic/networker/pkg/proxy"
	"github.com/spf13/cobra"
)

var (
	upStream, listenOn string
	proxyCmd           = &cobra.Command{
		Use:     "proxy",
		Aliases: []string{"p"},
		Example: "TODO: proxy cmd example",
		Short:   "forward network traffic and splice connections together",
		Run: func(cmd *cobra.Command, args []string) {
			if err := proxy.Run(listenOn, up); err != nil {
				log.Printf("failed to run proxy - err : %s\n", err)
				cmd.Usage()
				return
			}
		},
	}
)

func init() {
	proxyCmd.Flags().StringVarP(&upStream, "upstream", "u", upStream, "<host>:<port> to proxy traffic to")
	proxyCmd.Flags().StringVarP(&listenOn, "listen-on", "l", listenOn, "port for proxy to listen on")
	Networker.AddCommand(proxyCmd)
}
