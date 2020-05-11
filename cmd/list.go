package cmd

import (
	"log"

	"github.com/fuskovic/networker/pkg/list"
	"github.com/spf13/cobra"
)

var (
	device                          string
	myLocalIP, myRemoteIP, myRouter bool
	rootErr                         = "failed to find devices - are you root?"
	longListEx                      = ""
	listCmd                         = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list information on connected device(s).",
		Run: func(cmd *cobra.Command, args []string) {
			if myLocalIP {
				list.LocalIP()
			}

			if myRemoteIP {
				list.RemoteIP()
			}

			if myRouter {
				list.Router()
			}

			if device != "" {
				if err := list.Device(device); err != nil {
					log.Println(rootErr)
					cmd.Usage()
				}
			} else {
				if err := list.AllDevices(); err != nil {
					log.Println(rootErr)
					cmd.Usage()
				}
			}
		},
	}
)

func init() {
	listCmd.Flags().StringVarP(&device, "device", "d", "", "name of network interface device")
	listCmd.Flags().BoolVar(&myLocalIP, "my-local-ip", myLocalIP, "enable this to list the local IP address of this node")
	listCmd.Flags().BoolVar(&myRemoteIP, "my-remote-ip", myRemoteIP, "enable this to list the remote IP address of this node")
	listCmd.Flags().BoolVar(&myRouter, "my-router", myRouter, "enable this to list the IP address of the router for this subnet")
	Networker.AddCommand(listCmd)
}
