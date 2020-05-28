package request

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

const (
	jsonExt = ".json"
	xmlExt  = ".xml"
)

var supportedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

// Config collects the command parameters for the request sub-command.
type Config struct {
	URL, Method, File string
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
		buf         bytes.Buffer
		contentType string
		data        []byte
		err         error
	)

	if c.File == "" {
		return buf, nil
	}

	ext := path.Ext(c.File)

	switch ext {
	case jsonExt:
		contentType = "Content-Type:application/json"
	case xmlExt:
		contentType = "Content-Type:application/xml"
	default:
		return buf, fmt.Errorf("%s is an unsupported file format", ext)
	}

	data, err = ioutil.ReadFile(c.File)
	if err != nil {
		return buf, err
	}

	c.Headers = append(c.Headers, contentType)
	return *bytes.NewBuffer(data), nil
}
