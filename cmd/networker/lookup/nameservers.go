package lookup

import (
	"context"
	"fmt"
	"os"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"github.com/fuskovic/networker/internal/resolve"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"
)

type nameserversCmd struct {
	hostname string
}

func (cmd *nameserversCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "nameservers",
		Usage: "[flags]",
		Desc:  "Lookup nameservers for the provided hostname.",
	}
}

func (cmd *nameserversCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVar(&cmd.hostname, "hostname", "", "Hostname to lookup nameservers for.")
}

func (cmd *nameserversCmd) Run(fl *pflag.FlagSet) {
	if cmd.hostname == "" {
		fl.Usage()
		flog.Error("hostname not provided")
		return
	}

	nameservers, err := resolve.NameServersByHostName(cmd.hostname)
	if err != nil {
		fl.Usage()
		flog.Error("lookup failed: %v", err)
		return
	}

	var nsFields []slog.Field
	for i := range nameservers {
		nsFields = append(nsFields, slog.F(fmt.Sprintf("nameserver %d", i+1), nameservers[i].Host))
	}
	slog.Make(sloghuman.Sink(os.Stdout)).Info(context.Background(), "lookup successful", nsFields...)
}
