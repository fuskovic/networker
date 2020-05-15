package cmd

import (
	"log"

	"github.com/fuskovic/networker/pkg/list"
	"github.com/spf13/cobra"
)

var (
	device                               string
	myLocalIP, myRemoteIP, myRouter, all bool
	longListEx                           = ""
	listCmd                              = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list information on connected device(s).",
		Run: func(cmd *cobra.Command, args []string) {
			if myLocalIP {
				if err := list.LocalIP(); err != nil {
					log.Println("failed to get local IP", "err =", err)
				}
			}

			if myRemoteIP {
				if err := list.RemoteIP(); err != nil {
					log.Println("failed to get remote IP", "err =", err)
				}
			}

			if myRouter {
				if err := list.Router(); err != nil {
					log.Println("failed to get gateway IP", "err =", err)
				}
			}

			if device != "" {
				if err := list.Device(device); err != nil {
					log.Println("failed to find device", "device =", device, "err =", err)
					cmd.Usage()
				}
			}

			if all {
				if err := list.AllDevices(); err != nil {
					log.Println("failed to find devices", "err =", err)
					cmd.Usage()
				}
			}
		},
	}
)

func init() {
	listCmd.Flags().StringVarP(&device, "device", "d", "", "list details of a specific network interface device by name")
	listCmd.Flags().BoolVar(&myLocalIP, "my-local-ip", myLocalIP, "enable this to list the local IP address of this node")
	listCmd.Flags().BoolVar(&myRemoteIP, "my-remote-ip", myRemoteIP, "enable this to list the remote IP address of this node")
	listCmd.Flags().BoolVar(&myRouter, "my-router", myRouter, "enable this to list the IP address of the gateway for this network")
	listCmd.Flags().BoolVarP(&all, "all", "a", all, "enable this to list all connected network interface devices and associated information")
	Networker.AddCommand(listCmd)
}
