package cmd

import (
	"log"

	"github.com/fuskovic/networker/pkg/list"
	"github.com/spf13/cobra"
)

var (
	device  string
	me, all bool

	listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Example: listExample,
		Short:   "List information on connected network devices.",
		Run: func(cmd *cobra.Command, args []string) {
			if me {
				if err := list.Me(); err != nil {
					log.Println("failed to list this machine's information", "err =", err)
				}
			}

			if all {
				log.Println("pinging network devices...")
				if err := list.AllDevices(); err != nil {
					log.Println("failed to find devices", "err =", err)
					cmd.Usage()
				}
			}
		},
	}
)

func init() {
	listCmd.Flags().BoolVarP(&me, "me", "m", me, "List the name, local IP, remote IP, and router IP for this machine and the network it's connected to.")
	listCmd.Flags().BoolVarP(&all, "all", "a", all, "List the IP, hostname, and connection status of all devices on this network. (must be run as root)")
	Networker.AddCommand(listCmd)
}
