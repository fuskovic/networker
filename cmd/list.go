package cmd

import (
	"fmt"

	"github.com/fuskovic/networker/pkg/list"
	"github.com/spf13/cobra"
)

var (
	listCfg = &list.Config{}

	listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Example: listExample,
		Short:   "List information on connected network devices.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := list.Run(listCfg); err != nil {
				fmt.Println(err)
				cmd.Usage()
			}
		},
	}
)

func init() {
	listCmd.Flags().BoolVarP(&listCfg.Me, "me", "m", listCfg.Me, "List the name, local IP, remote IP, and router IP for this machine and the network it's connected to.")
	listCmd.Flags().BoolVarP(&listCfg.All, "all", "a", listCfg.All, "List the IP, hostname, and connection status of all devices on this network. (must be run as root)")
	Networker.AddCommand(listCmd)
}
