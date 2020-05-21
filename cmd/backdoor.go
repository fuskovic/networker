package cmd

import (
	"fmt"

	"github.com/fuskovic/networker/pkg/backdoor"
	"github.com/spf13/cobra"
)

var (
	backDoorCfg = &backdoor.Config{}

	backDoorCmd = &cobra.Command{
		Use:     "backdoor",
		Aliases: []string{"bd", "b"},
		Example: backDoorExample,
		Short:   "Create and connect to backdoors to gain shell access over TCP.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := backdoor.Run(backDoorCfg); err != nil {
				fmt.Println(err)
				cmd.Usage()
			}
		},
	}
)

func init() {
	backDoorCmd.Flags().BoolVar(&backDoorCfg.Create, "create", backDoorCfg.Create, "Create a TCP backdoor. (must be used with the --port flag)")
	backDoorCmd.Flags().BoolVar(&backDoorCfg.Connect, "connect", backDoorCfg.Connect, "Connect to a TCP backdoor. (must be used with the --address flag)")
	backDoorCmd.Flags().StringVarP(&backDoorCfg.Address, "address", "a", backDoorCfg.Address, "Address of a remote target to connect to. (format: <host>:<port>)")
	backDoorCmd.Flags().IntVarP(&backDoorCfg.Port, "port", "p", backDoorCfg.Port, "Port number to listen for connections on. (format: 80)")
	switch {
	case backDoorCfg.Create:
		backDoorCmd.MarkFlagRequired("port")
	case backDoorCfg.Connect:
		backDoorCmd.MarkFlagRequired("address")
	}
	Networker.AddCommand(backDoorCmd)
}
