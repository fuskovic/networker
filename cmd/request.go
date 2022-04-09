package cmd

import (
	"fmt"
	"io/ioutil"

	"net/http"
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
	silence   bool
	verbose   bool
)

func init() {
	requestCmd.Flags().StringSliceVarP(&headers, "headers", "H", headers, "Request headers.(format: key:value,key:value,key:value)")
	requestCmd.Flags().StringVarP(&method, "method", "m", "GET", "Request method.")
	requestCmd.Flags().StringVarP(&body, "body", "b", body, "Request body. (you can use a JSON string literal or a path to a json file)")
	requestCmd.Flags().StringVarP(&filePaths, "files", "f", filePaths, "Upload form file(s). (format: formname=path/to/file1,path/to/file2,path/to/file3)")
	requestCmd.Flags().BoolVarP(&silence, "silence", "s", silence, "Omit output of everything except JSON response body.")
	requestCmd.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "Include response headers in output.")
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
			-b '{"field": "doesntmatter"}' \
			https://some-url.com/api/v1/some/endpoint

	POST request using json body sourced from file:
		networker request \
			-H "Authorization: Bearer doesntmatter" \
			-m post -b /path/to/file.json \
			https://some-url.com/api/v1/some/endpoint

	PUT request for file upload:
		networker request \
			-H "Authorization: Bearer doesntmatter" \
			-m put \
			-f formname=/path/to/file1.jpeg \
			https://some-url.com/api/v1/some/endpoint

	PUT request for uploading multiple files:
		networker request \
			-H "Authorization: Bearer doesntmatter" \
			-m put \
			-f formname=/path/to/file1.jpeg,/path/to/file2.jpeg,/path/to/file3.jpeg \
			https://some-url.com/api/v1/some/endpoint
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		req, err := request.New(
			&request.Config{
				Headers:   headers,
				URL:       args[0],
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

		if !silence {
			fmt.Printf("status_code: %d\n", resp.StatusCode)
			fmt.Printf("latency(ms): %s\n", time.Since(started))
		}

		if verbose {
			for k, v := range resp.Header {
				fmt.Printf("%s: %s\n", k, v)
			}
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			usage.Fatalf(cmd, "failed to read response body : %v", err)
		}
		println(string(data))
	},
}
