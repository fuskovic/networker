package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"

	"github.com/fuskovic/networker/internal/request"
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
	req, err := request.New(cmd.cfg())
	if err != nil {
		fl.Usage()
		flog.Error("failed to build request : %v", err)
		return
	}

	client := http.Client{Timeout: time.Duration(cmd.timeOut) * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		fl.Usage()
		flog.Error("failed to send HTTP request : %v", err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		flog.Error("failed to read response body : %v", err)
		fl.Usage()
		return
	}

	sloghuman.Make(os.Stdout).Info(context.Background(), "received server response",
		slog.F("method", req.Method),
		slog.F("status", fmt.Sprintf("%d - %s", resp.StatusCode, http.StatusText(resp.StatusCode))),
		slog.F("response", string(data)),
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
