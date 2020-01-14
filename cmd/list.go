package cmd

import (
	"fmt"
	"log"

	p "github.com/google/gopacket/pcap"
	"github.com/spf13/cobra"
)

var (
	device     string
	allDevices []p.Interface

	list = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Long:    "list information of device(s)",
		Run: func(cmd *cobra.Command, args []string) {
			devices, err := p.FindAllDevs()
			if err != nil {
				fmt.Println("failed to find devices - are you root?")
				log.Fatalf("error : %v\n", err)
			}

			allDevices = devices
			if device != "" {
				listDevice(device)
			} else {
				listAllDevices()
			}
		},
	}
)

func init() {
	list.Flags().StringVarP(&device, "device", "d", "", "name of network interface device")
	Networker.AddCommand(list)
}

func listDevice(name string) {
	for _, d := range allDevices {
		if d.Name == name {
			print(d)
		}
	}
}

func listAllDevices() {
	for _, d := range allDevices {
		print(d)
	}
}

func print(d p.Interface) {
	fmt.Printf("\nName: %s\nDescription: %s\n", d.Name, d.Description)
	for _, a := range d.Addresses {
		fmt.Printf("\n- IP address: %s\n- Subnet mask: %s\n", a.IP, a.Netmask)
	}
}
