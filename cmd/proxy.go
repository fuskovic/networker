package cmd

import (
	"fmt"

	"github.com/fuskovic/networker/pkg/proxy"
	"github.com/spf13/cobra"
)

var (
	upStream     string
	listenOn     int
	proxyLong    = "\nnetworker proxy --listen-on <port> -upstream <host>:<port>\n"
	proxyShort   = "\nnetworker p -l <port> -u <host>:<port>"
	format       = "\nlong format:\n%s\nshort format:\n%s"
	proxyExmaple = fmt.Sprintf(format, proxyLong, proxyShort)
	proxyCmd     = &cobra.Command{
		Use:     "proxy",
		Aliases: []string{"p"},
		Example: proxyExmaple,
		Short:   "forward network traffic from one network connection to another",
		Run: func(cmd *cobra.Command, args []string) {
			proxy.Run(listenOn, upStream)
		},
	}
)

func init() {
	proxyCmd.Flags().StringVarP(&upStream, "upstream", "u", upStream, "<host>:<port> to proxy traffic to")
	proxyCmd.Flags().IntVarP(&listenOn, "listen-on", "l", listenOn, "port for proxy to listen on")
	proxyCmd.MarkFlagRequired("upstream")
	proxyCmd.MarkFlagRequired("listen-on")
	Networker.AddCommand(proxyCmd)
}
