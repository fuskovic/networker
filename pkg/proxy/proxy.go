package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const tcp = "tcp"

// Run initializes and starts a new proxy server and forwards traffic from the listener to the upstream server.
func Run(listenOn int, upStream string) {
	port := fmt.Sprintf(":%d", listenOn)

	fmt.Printf("starting listener on %s...\n", port)
	lsnr, err := net.Listen(tcp, port)
	if err != nil {
		log.Println(err)
		return
	}
	defer lsnr.Close()

	log.Printf("dialing %s\n", upStream)

	upStr, err := net.Dial(tcp, upStream)
	if err != nil {
		log.Println(err)
		return
	}
	defer upStr.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		for {
			sig := <-c
			log.Printf("received termination signal : %v\n", sig)
			lsnr.Close()
		}
	}()

	log.Println("proxy started")
	for {
		conn, err := lsnr.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()

		go io.Copy(conn, upStr)
		io.Copy(upStr, conn)
	}
}
