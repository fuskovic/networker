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

var signals = []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}

// Config collects the command parameters for the proxy sub-command.
type Config struct {
	ListenOn int
	UpStream string
}

// Run initializes and starts a new proxy server and forwards traffic from the listener to the upstream server.
func Run(cfg *Config) error {
	port := fmt.Sprintf(":%d", cfg.ListenOn)

	log.Printf("starting listener on %s...\n", port)

	lsnr, err := net.Listen(tcp, port)
	if err != nil {
		return err
	}
	defer lsnr.Close()

	log.Printf("dialing %s\n", cfg.UpStream)

	upStr, err := net.Dial(tcp, cfg.UpStream)
	if err != nil {
		return err
	}
	defer upStr.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)

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
			return err
		}
		defer conn.Close()

		go io.Copy(conn, upStr)
		io.Copy(upStr, conn)
	}
}
