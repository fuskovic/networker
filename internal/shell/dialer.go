package shell

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// Dial connects to the shell that it expects to be served at addr.
func Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to dial %q : %w", addr, err)
	}
	defer conn.Close()

	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Printf("failed to get previous state of terminal: %s", err)
	}

	defer func() {
		if oldState != nil {
			_ = terminal.Restore(int(os.Stdin.Fd()), oldState)
		}
	}()

	errChan := make(chan error, 1)

	go func() {
		if _, err := io.Copy(os.Stdout, conn); err != nil && !errors.Is(err, net.ErrClosed) {
			errChan <- fmt.Errorf("failed to read output from connection: %s\n", err)
			return
		}
		errChan <- nil
	}()

	go func() {
		if _, err = io.Copy(conn, os.Stdin); err != nil && !errors.Is(err, syscall.EPIPE) {
			errChan <- fmt.Errorf("failed to write input to connection: %w", err)
			return
		}
		errChan <- nil
	}()

	return <-errChan
}
