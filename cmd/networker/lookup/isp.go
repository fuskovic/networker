package lookup

import (
	"context"
	"os"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"

	"github.com/fuskovic/networker/internal/resolve"
)

type ispCmd struct {
	host string
}

func (cmd *ispCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "isp",
		Usage: "[flags]",
		Desc:  "Lookup the internet service provider of a remote host.",
	}
}

func (cmd *ispCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVar(&cmd.host, "host", "", "IP address or hostname to get the network address for.")
}

func (cmd *ispCmd) Run(fl *pflag.FlagSet) {
	if cmd.host == "" {
		fl.Usage()
		flog.Error("no host provided")
		return
	}

	_, ip, err := resolve.HostAndAddr(cmd.host)
	if err != nil {
		fl.Usage()
		flog.Error("%q is an invalid host: %w", cmd.host, err)
		return
	}

	if resolve.IsPrivate(ip) {
		fl.Usage()
		flog.Error("%q is not a remote ip", ip)
		return
	}

	isp, err := resolve.ServiceProvider(ip)
	if err != nil {
		fl.Usage()
		flog.Error("failed to resolve internet service provider for %q: %w", cmd.host, err)
		return
	}

	var (
		ctx = context.Background()
		log = slog.Make(sloghuman.Sink(os.Stdout))
	)

	for _, field := range isp.Fields() {
		log.Info(ctx, "isp-lookup", field)
	}
}
