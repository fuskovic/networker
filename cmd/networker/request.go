package networker

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"github.com/spf13/pflag"
	"go.coder.com/cli"

	"github.com/fuskovic/networker/internal/request"
	"github.com/fuskovic/networker/internal/usage"
)

type requestCmd struct {
	url     string
	method  string
	file    string
	headers []string
	timeOut int
}

func (cmd *requestCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "request",
		Usage:   "[flags]",
		Aliases: []string{"r", "req"},
		Desc:    "Send an HTTP request.",
	}
}

func (cmd *requestCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringSliceVarP(&cmd.headers, "add-headers", "a", cmd.headers, "Add a list of comma-separated request headers. (format : key:value,key:value,etc...)")
	fl.StringVarP(&cmd.url, "url", "u", cmd.url, "URL to send request.")
	fl.StringVarP(&cmd.method, "method", "m", "GET", "Specify method. (supported methods include GET, POST, PUT, PATCH, and DELETE)")
	fl.StringVarP(&cmd.file, "file", "f", cmd.file, "Path to JSON or XML file to use for request body. (content-type headers for each file-type are set automatically)")
	fl.IntVarP(&cmd.timeOut, "time-out", "t", 3, "Specify number of seconds for time-out.")
}

func (cmd *requestCmd) Run(fl *pflag.FlagSet) {
	started := time.Now()

	req, err := request.New(cmd.cfg())
	if err != nil {
		usage.Fatalf(fl, "failed to build request : %v", err)
	}

	client := http.Client{Timeout: time.Duration(cmd.timeOut) * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		usage.Fatalf(fl, "failed to send HTTP request : %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		usage.Fatalf(fl, "failed to read response body : %v", err)
	}

	var successful bool
	if resp.StatusCode < 400 {
		successful = true
	}

	ctx := context.Background()
	log := slog.Make(sloghuman.Sink(os.Stdout))

	log.Info(ctx, "response",
		slog.F("successful", successful),
		slog.F("method", req.Method),
		slog.F("status", resp.Status),
		slog.F("elapsed-time", time.Since(started)),
		slog.F("body", string(data)),
	)
}

func (cmd *requestCmd) cfg() *request.Cfg {
	return &request.Cfg{
		Headers: cmd.headers,
		URL:     cmd.url,
		Method:  cmd.method,
		File:    cmd.file,
		Seconds: cmd.timeOut,
	}
}
