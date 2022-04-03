package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/fuskovic/networker/internal/request"
	"github.com/fuskovic/networker/internal/usage"
)

var (
	method    string
	body      string
	filePaths string
	headers   []string
	jsonOnly  bool
)

func init() {
	requestCmd.Flags().StringSliceVarP(&headers, "headers", "H", headers, "Request headers.(format: key:value,key:value,key:value)")
	requestCmd.Flags().StringVarP(&method, "method", "m", "GET", "Request method.")
	requestCmd.Flags().StringVarP(&body, "body", "b", body, "Request body. (you can use a JSON string literal or a path to a json file)")
	requestCmd.Flags().StringVarP(&filePaths, "files", "f", filePaths, "Files to upload. (format: formname=path/to/file1,path/to/file2,path/to/file3)")
	requestCmd.Flags().BoolVarP(&jsonOnly, "json-only", "j", jsonOnly, "Only output json response body.")
	Root.AddCommand(requestCmd)
}

var requestCmd = &cobra.Command{
	Use:     "request",
	Aliases: []string{"r", "req"},
	Short:   "Send an HTTP request.",
	Example: `
POST request using json body sourced from stdin:
	networker request \
		-H "Authorization: Bearer doesntmatter" \
		-m post \
		-b '{"field": "doesntmatter"}'

POST request using json body sourced from file:
	networker request \
		-H "Authorization: Bearer doesntmatter" \
		-m post -b /path/to/file.json

PUT request for file upload:
	networker request \
		-H "Authorization: Bearer doesntmatter" \
		-m put \
		-f=/path/to/file1.jpeg

PUT request for uploading multiple files:
	networker request \
		-H "Authorization: Bearer doesntmatter" \
		-m put \
		-f=/path/to/file1.jpeg,/path/to/file2.jpeg,/path/to/file3.jpeg
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		req, err := request.NewNetworkerCraftedHTTPRequest(
			&request.Config{
				Headers:   headers,
				URL:       os.Args[len(os.Args)-1],
				Method:    method,
				Body:      body,
				FilePaths: filePaths,
			},
		)

		if err != nil {
			usage.Fatalf(cmd, "failed to build request : %v", err)
		}

		started := time.Now()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			usage.Fatalf(cmd, "failed to send HTTP request : %v", err)
		}
		defer resp.Body.Close()

		ended := time.Since(started)
		if !jsonOnly {
			log.Printf("received response in: %s\nstatus: %s\n", ended, resp.Status)
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			usage.Fatalf(cmd, "failed to read response body : %v", err)
		}

		b := bytes.NewBuffer(nil)
		if err := json.Indent(b, data, "", " "); err != nil {
			// it's possible that a non-json response is received so if we fail to encode here just output whatever came back
			b = bytes.NewBuffer(data)
		}

		if jsonOnly {
			println(b.String())
		} else {
			log.Printf("response:\n%s\n", b.String())
		}
	},
}
