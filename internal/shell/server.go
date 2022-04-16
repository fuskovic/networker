package shell

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/creack/pty"
)

// Serve serves a shell on the designated port.
func Serve(shell, port string) error {
	if !isSupportedShell(shell) {
		return fmt.Errorf("shell %q is not supported", shell)
	}

	if _, err := exec.LookPath(shell); err != nil {
		return fmt.Errorf("shell %q does not exist on system: %w", shell, err)
	}

	if !isValidPort(port) {
		return fmt.Errorf("%q is not a valid port", port)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	errChan := make(chan error, 1)

	go func() {
		<-c
		println()
		log.Println("received interrupt signal")
		log.Println("shutting down")
		errChan <- nil
	}()

	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}
	defer l.Close()

	go func() {
		for {
			proc, err := pty.Start(exec.Command(shell))
			if err != nil {
				errChan <- fmt.Errorf("failed to start new shell process: %w", err)
				return
			}
			defer proc.Close()

			ch := make(chan os.Signal, 1)
			signal.Notify(ch, syscall.SIGWINCH)

			go func() {
				for range ch {
					if err := pty.InheritSize(os.Stdin, proc); err != nil {
						log.Printf("error resizing: %s", err)
					}
				}
			}()

			// initial resize
			ch <- syscall.SIGWINCH
			defer func() {
				signal.Stop(ch)
				close(ch)
			}()

			log.Printf("serving a new shell process on localhost:%s\n", port)

			conn, err := l.Accept()
			if err != nil {
				errChan <- fmt.Errorf("failed to accept incoming connection: %w", err)
				return
			}
			defer conn.Close()

			go handleConnection(conn, proc)
		}
	}()
	return <-errChan
}

func handleConnection(conn net.Conn, proc *os.File) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close connection for %s: %s\n", conn.RemoteAddr(), err)
		}

		if err := proc.Close(); err != nil {
			log.Printf("failed to kill process started by %s: %s\n", conn.RemoteAddr(), err)
		}
	}()

	connectedAt := time.Now().UTC()
	log.Printf("client %s connected at: %s\n", conn.RemoteAddr(), connectedAt)

	go func() {
		if _, err := io.Copy(proc, conn); err != nil && !errors.Is(err, net.ErrClosed) {
			log.Printf("failed to read from client connection: %+v\n", err)
		}
	}()

	if _, err := io.Copy(conn, proc); err != nil {
		log.Printf("failed to write output to connection: %s", err)
	}

	log.Printf("client %s disconnected after %s\n", conn.RemoteAddr(), time.Since(connectedAt))
}

func isSupportedShell(targetShell string) bool {
	for _, sh := range []string{"bash", "zsh", "sh", "fish"} {
		if sh == targetShell {
			return true
		}
	}
	return false
}

func isValidPort(port string) bool {
	p, err := strconv.Atoi(port)
	if err != nil {
		return false
	}
	return p > -1 && p < 65536
}