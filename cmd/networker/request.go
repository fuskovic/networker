package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"

	req "github.com/fuskovic/networker/internal/request"
)

type requestCmd struct {
	url     string
	method  string
	file    string
	headers []string
	timeOut int
}

// Spec returns a command spec containing a description of it's usage.
func (c *requestCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "request",
		Usage: "[flags]",
		Desc:  "Send an HTTP request.",
	}
}

// RegisterFlags initializes how a flag set is processed for a particular command.
func (c *requestCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringSliceVarP(&c.headers, "add-headers", "a", c.headers, "Add a list of comma-separated request headers. (format : key:value,key:value,etc...)")
	fl.StringVarP(&c.url, "url", "u", c.url, "URL to send request.")
	fl.StringVarP(&c.method, "method", "m", "GET", "Specify method. (supported methods include GET, POST, PUT, PATCH, and DELETE)")
	fl.StringVarP(&c.file, "file", "f", c.file, "Path to JSON or XML file to use for request body. (content-type headers for each file-type are set automatically)")
	fl.IntVarP(&c.timeOut, "time-out", "t", 3, "Specify number of seconds for time-out.")
}

// Run crafts an HTTP request out of the specified flag set, sends it, and outputs the response.
func (c *requestCmd) Run(fl *pflag.FlagSet) {
	flog.Info("building request")

	r, err := req.New(c.cfg())
	if err != nil {
		flog.Error("failed to build request : %v", err)
		fl.Usage()
		return
	}

	flog.Info("sending %s request", r.Method)

	t := time.Duration(c.timeOut) * time.Second
	client := http.Client{Timeout: t}

	resp, err := client.Do(r)
	if err != nil {
		flog.Error("failed to send HTTP request : %v", err)
		fl.Usage()
		return
	}
	defer resp.Body.Close()

	flog.Info("checking response\n\n")

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		flog.Error("failed to read body : %v", err)
		fl.Usage()
		return
	}

	fmt.Fprint(os.Stdout, string(data)+"\n\n")

	msg := http.StatusText(resp.StatusCode)
	format := "%d: %s"

	if resp.StatusCode == http.StatusOK {
		flog.Success(format, resp.StatusCode, msg)
	} else {
		flog.Error(format, resp.StatusCode, msg)
	}
}

func (c *requestCmd) cfg() *req.Cfg {
	return &req.Cfg{
		Headers: c.headers,
		URL:     c.url,
		Method:  c.method,
		File:    c.file,
		Seconds: c.timeOut,
	}
}
