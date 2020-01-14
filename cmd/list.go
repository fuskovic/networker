package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/fuskovic/networker/pkg/list"
)

var (
	device     string
	rootErr = "failed to find devices - are you root?"
	listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:    "list information on connected device(s).",
		Run: func(cmd *cobra.Command, args []string) {
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
	Networker.AddCommand(listCmd)
}