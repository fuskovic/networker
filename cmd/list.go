package cmd

import (
	"log"

	"github.com/fuskovic/networker/pkg/list"
	"github.com/spf13/cobra"
)

var (
	device      string
	me, all     bool
	listExample = "\nnetworker ls --me -a"
	listCmd     = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Example: listExample,
		Short:   "list information on connected device(s).",
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
	listCmd.Flags().BoolVar(&me, "me", me, "enable this to list the name, local IP, remote IP, and router IP for this machine")
	listCmd.Flags().BoolVarP(&all, "all", "a", all, "enable this to list all connected network interface devices and associated information")
	Networker.AddCommand(listCmd)
}
