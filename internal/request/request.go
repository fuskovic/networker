package request

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

const (
	// supported file formats for provisioning the request payload from a file.
	jsonExt = ".json"
	xmlExt  = ".xml"
)

// Cfg is used to determine how we build an HTTP request.
type Cfg struct {
	Headers []string
	URL     string
	Method  string
	File    string
	Seconds int
}

// New validates the config and builds a new HTTP request.
// This includes provisioning the request with the headers and payload provided from the config.
func New(cfg *Cfg) (*http.Request, error) {
	if err := cfg.valid(); err != nil {
		return nil, fmt.Errorf("invalid request : %v", err)
	}

	payload := new(bytes.Buffer)
	if requiresPayload(cfg.Method) {
		body, err := buildPayload(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to build request payload : %v", err)
		}
		payload = body
	}

	req, err := http.NewRequest(cfg.Method, cfg.URL, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to craft request : %v", err)
	}

	if err := addHeaders(cfg, req); err != nil {
		return nil, fmt.Errorf("failed to add headers: %w", err)
	}
	return req, nil
}

func (cfg *Cfg) valid() error {
	switch {
	case cfg.URL == "":
		return errors.New("no url endpoint specified")
	case !isSupported(cfg.Method):
		return fmt.Errorf("%s is an invalid request method", cfg.Method)
	default:
		if !hasProtoScheme(cfg.URL) {
			cfg.URL = "https://" + cfg.URL
		}
	}
	return nil
}

func buildPayload(cfg *Cfg) (*bytes.Buffer, error) {
	var contentType string
	if cfg.File == "" {
		return new(bytes.Buffer), nil
	}

	switch ext := path.Ext(cfg.File); ext {
	case jsonExt:
		contentType = "Content-Type:application/json"
	case xmlExt:
		contentType = "Content-Type:application/xml"
	default:
		return new(bytes.Buffer), fmt.Errorf("%q unsupported file format ", ext)
	}

	cfg.Headers = append(cfg.Headers, contentType)

	b, err := ioutil.ReadFile(cfg.File)
	if err != nil {
		return new(bytes.Buffer), fmt.Errorf("failed to read file %q: %w", cfg.File, err)
	}
	return bytes.NewBuffer(b), nil
}

func addHeaders(cfg *Cfg, r *http.Request) error {
	for _, h := range cfg.Headers {
		args := strings.Split(h, ":")
		if (len(args) % 2) != 0 {
			return errors.New("uneven number of key/value pairs")
		}
		r.Header.Set(args[0], args[1])
	}
	return nil
}

func isSupported(method string) bool {
	for _, m := range []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	} {
		if m == method {
			return true
		}
	}
	return false
}

func requiresPayload(method string) bool {
	return method == http.MethodPatch || method == http.MethodPut || method == http.MethodPost
}

func hasProtoScheme(url string) bool {
	return strings.Contains(url, "http://") || strings.Contains(url, "https://")
}
