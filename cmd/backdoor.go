package cmd

import (
	"github.com/fuskovic/networker/pkg/backdoor"
	"github.com/spf13/cobra"
)

var (
	create, connect bool
	port            int
	address         string
	backDoorCmd     = &cobra.Command{
		Use:     "backdoor",
		Aliases: []string{"bd", "b"},
		Example: "TODO: add backdoor cmd syntax",
		Short:   "Create and connect to backdoors over TCP",
		Run: func(cmd *cobra.Command, args []string) {
			switch {
			case create:
				backdoor.Create(port)
			case connect:
				// TODO : implement connect logic
			}
		},
	}
)

func init() {
	backDoorCmd.Flags().BoolVar(&create, "create", create, "create a TCP backdoor(must be used with --port flag)")
	backDoorCmd.Flags().BoolVar(&connect, "connect", connect, "connect to a TCP backdoor(must be used with --address flag)")
	backDoorCmd.Flags().StringVarP(&address, "address", "a", address, "address of remote target to connect to(format: <host>:<port>))")
	backDoorCmd.Flags().IntVarP(&port, "port", "p", port, "port number to listen for connections on")
	switch {
	case create:
		backDoorCmd.MarkFlagRequired("port")
	case connect:
		backDoorCmd.MarkFlagRequired("address")
	}
	Networker.AddCommand(backDoorCmd)
}
