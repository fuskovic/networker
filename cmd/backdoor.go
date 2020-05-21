package cmd

import (
	"github.com/fuskovic/networker/pkg/backdoor"
	"github.com/spf13/cobra"
)

var (
	create, connect bool
	port            int
	address         string

	backDoorCmd = &cobra.Command{
		Use:     "backdoor",
		Aliases: []string{"bd", "b"},
		Example: backDoorExample,
		Short:   "Create and connect to backdoors to gain shell access over TCP.",
		Run: func(cmd *cobra.Command, args []string) {
			switch {
			case create:
				backdoor.Create(port)
			case connect:
				backdoor.Connect(address)
			}
		},
	}
)

func init() {
	backDoorCmd.Flags().BoolVar(&create, "create", create, "Create a TCP backdoor. (must be used with the --port flag)")
	backDoorCmd.Flags().BoolVar(&connect, "connect", connect, "Connect to a TCP backdoor. (must be used with the --address flag)")
	backDoorCmd.Flags().StringVarP(&address, "address", "a", address, "Address of a remote target to connect to. (format: <host>:<port>)")
	backDoorCmd.Flags().IntVarP(&port, "port", "p", port, "Port number to listen for connections on. (format: 80)")
	switch {
	case create:
		backDoorCmd.MarkFlagRequired("port")
	case connect:
		backDoorCmd.MarkFlagRequired("address")
	}
	Networker.AddCommand(backDoorCmd)
}
