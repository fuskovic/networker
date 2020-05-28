package request

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"
)

const jsonExt = ".json"

var supportedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

// Config collects the command parameters for the request sub-command.
type Config struct {
	URL, Method, Data string
	Headers           []string
	TimeOut           int
}

func (c *Config) validMethod() bool {
	for _, m := range supportedMethods {
		if m == c.Method {
			return true
		}
	}
	return false
}

func (c *Config) hasProtoScheme() bool {
	has := func(s string) bool { return strings.Contains(c.URL, s) }
	return has("http://") || has("https://")
}

func (c *Config) buildBody() (bytes.Buffer, error) {
	var (
		buf  bytes.Buffer
		data []byte
		err  error
	)

	if c.Data == "" {
		return buf, nil
	}

	if path.Ext(c.Data) != jsonExt {
		data = []byte(c.Data)
	} else {
		data, err = ioutil.ReadFile(c.Data)
		if err != nil {
			return buf, err
		}
	}

	enc := json.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return buf, err
	}
	return buf, nil
}
