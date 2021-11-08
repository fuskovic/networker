package networker

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/pflag"
	"go.coder.com/cli"

	"github.com/fuskovic/networker/internal/request"
	"github.com/fuskovic/networker/internal/usage"
)

type requestCmd struct {
	method        string
	body          string
	multiPartForm string
	headers       []string
	jsonOnly      bool
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
	fl.StringSliceVarP(&cmd.headers, "headers", "H", cmd.headers, "Request headers.(format(no quotes): key:value,key:value,key:value)")
	fl.StringVarP(&cmd.method, "method", "m", "GET", "Request method.")
	fl.StringVarP(&cmd.body, "body", "b", cmd.body, "Request body. (you can use a JSON string literal or a path to a json file)")
	fl.StringVarP(&cmd.multiPartForm, "upload", "u", cmd.multiPartForm, "Multi-part form. (format: formname=path/to/file1,path/to/file2,path/to/file3)")
	fl.BoolVarP(&cmd.jsonOnly, "json-only", "j", cmd.jsonOnly, "Only output json.")
}

func (cmd *requestCmd) Run(fl *pflag.FlagSet) {
	if len(os.Args) < 3 {
		usage.Fatal(fl, "url not provided")
	}

	req, err := request.NewNetworkerCraftedHTTPRequest(
		&request.Config{
			Headers:       cmd.headers,
			URL:           os.Args[len(os.Args)-1],
			Method:        cmd.method,
			Body:          cmd.body,
			MultiPartForm: cmd.multiPartForm,
		},
	)

	if err != nil {
		usage.Fatalf(fl, "failed to build request : %v", err)
	}

	client := http.DefaultClient
	started := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		usage.Fatalf(fl, "failed to send HTTP request : %v", err)
	}
	defer resp.Body.Close()

	ended := time.Since(started)
	if !cmd.jsonOnly {
		log.Printf("received response in: %s\nstatus: %s\n", ended, resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		usage.Fatalf(fl, "failed to read response body : %v", err)
	}

	b := bytes.NewBuffer(nil)
	if err := json.Indent(b, data, "", " "); err != nil {
		// it's possible that a non-json response is received so if we fail to encode here just output whatever came back
		b = bytes.NewBuffer(data)
	}

	if cmd.jsonOnly {
		println(b.String())
	} else {
		log.Printf("response:\n%s\n", b.String())
	}
}
