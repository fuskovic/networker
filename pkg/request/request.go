package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

const jsonExt = ".json"

// Config collects the command parameters for the request sub-command.
type Config struct{ URL, Method, Data string }

func (c *Config) validMethod() bool {
	return c.Method == "GET" || c.Method == "POST"
}

func (c *Config) hasProtoScheme() bool {
	has := func(s string) bool { return strings.Contains(c.URL, s) }
	return has("http://") || has("https://")
}

func (c *Config) buildBody() (*bytes.Buffer, error) {
	if c.Data == "" {
		return nil, nil
	}

	var (
		buf  *bytes.Buffer
		data []byte
		err  error
	)

	if path.Ext(c.Data) != jsonExt {
		data = []byte(c.Data)
	} else {
		data, err = ioutil.ReadFile(c.Data)
		if err != nil {
			return nil, err
		}
	}

	enc := json.NewEncoder(buf)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}

	return buf, nil
}

// Run executes the command logic for the request package.
func Run(cfg *Config) error {
	var resp *http.Response

	body, err := cfg.buildBody()
	if err != nil {
		return err
	}

	if !cfg.hasProtoScheme() {
		cfg.URL = "https://" + cfg.URL
	}

	switch cfg.Method {
	case "GET":
		resp, err = http.Get(cfg.URL)
	case "POST":
		resp, err = http.Post(cfg.URL, "application/json", body)
	default:
		return fmt.Errorf("%s is an unsupported request method", cfg.Method)
	}

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
	return nil
}
