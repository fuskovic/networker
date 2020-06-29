package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"
)

const (
	jsonExt = ".json"
	xmlExt  = ".xml"
)

var supportedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

type requestCmd struct {
	url, method, file string
	headers           []string
	timeOut           int
}

// Spec returns a command spec containing a description of it's usage.
func (cmd *requestCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "request",
		Usage: "[flags]",
		Desc:  "Send an HTTP request.",
	}
}

// RegisterFlags initializes how a flag set is processed for a particular command.
func (cmd *requestCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringSliceVarP(&cmd.headers, "add-headers", "a", cmd.headers, "Add a list of comma-separated request headers. (format : key:value,key:value,etc...)")
	fl.StringVarP(&cmd.url, "url", "u", cmd.url, "URL to send request.")
	fl.StringVarP(&cmd.method, "method", "m", "GET", "Specify method. (supported methods include GET, POST, PUT, PATCH, and DELETE)")
	fl.StringVarP(&cmd.file, "file", "f", cmd.file, "Path to JSON or XML file to use for request body. (content-type headers for each file-type are set automatically)")
	fl.IntVarP(&cmd.timeOut, "time-out", "t", 3, "Specify number of seconds for time-out.")
}

// Run crafts an HTTP request out of the specified flag set, sends it, and outputs the response.
func (cmd *requestCmd) Run(fl *pflag.FlagSet) {
	seconds := time.Duration(cmd.timeOut)
	timeOut := time.Duration(seconds * time.Second)
	client := http.Client{Timeout: timeOut}

	flog.Info("validating URL")

	if cmd.url == "" {
		flog.Error("No endpoint")
		fl.Usage()
		return
	}

	body, err := cmd.buildBody()
	if err != nil {
		flog.Error(err.Error())
		fl.Usage()
		return
	}

	if !cmd.hasProtoScheme() {
		cmd.url = "https://" + cmd.url
	}

	flog.Info("validating Method")

	if !cmd.validMethod() {
		flog.Error(fmt.Sprintf("%s is an invalid request method", cmd.method))
		fl.Usage()
		return
	}

	flog.Info("building request")

	req, err := http.NewRequest(cmd.method, cmd.url, &body)
	if err != nil {
		flog.Error(err.Error())
		fl.Usage()
		return
	}

	flog.Info("Adding headers")

	for _, h := range cmd.headers {
		flog.Info("%s", h)
		kvPair := strings.Split(h, ":")
		req.Header.Set(kvPair[0], kvPair[1])
	}

	flog.Info("Sending request")

	resp, err := client.Do(req)
	if err != nil {
		flog.Error(err.Error())
		fl.Usage()
		return
	}
	defer resp.Body.Close()

	flog.Info("Received response")

	msg := http.StatusText(resp.StatusCode)

	if resp.StatusCode >= 400 {
		flog.Error("%d:%s", resp.StatusCode, msg)
	} else {
		flog.Success("%d:%s", resp.StatusCode, msg)
	}

	io.Copy(os.Stdout, resp.Body)
	println()
}

func (cmd *requestCmd) validMethod() bool {
	for _, m := range supportedMethods {
		if m == cmd.method {
			return true
		}
	}
	return false
}

func (cmd *requestCmd) hasProtoScheme() bool {
	has := func(s string) bool { return strings.Contains(cmd.url, s) }
	return has("http://") || has("https://")
}

func (cmd *requestCmd) buildBody() (bytes.Buffer, error) {
	var (
		buf         bytes.Buffer
		contentType string
		data        []byte
		err         error
	)

	if cmd.file == "" {
		return buf, nil
	}

	ext := path.Ext(cmd.file)

	switch ext {
	case jsonExt:
		contentType = "Content-Type:application/json"
	case xmlExt:
		contentType = "Content-Type:application/xml"
	default:
		return buf, fmt.Errorf("%s is an unsupported file format", ext)
	}

	data, err = ioutil.ReadFile(cmd.file)
	if err != nil {
		return buf, err
	}

	cmd.headers = append(cmd.headers, contentType)
	return *bytes.NewBuffer(data), nil
}
