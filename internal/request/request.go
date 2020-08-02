package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	u "github.com/fuskovic/networker/internal/utils"
	"go.coder.com/flog"
)

const (
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
func New(cfg *Cfg) (*http.Request, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("failed validation : %v", err)
	}

	b, err := cfg.body()
	if err != nil {
		return nil, fmt.Errorf("failed to build body : %v", err)
	}

	req, err := http.NewRequest(cfg.Method, cfg.URL, &b)
	if err != nil {
		return nil, fmt.Errorf("failed to create new HTTP request : %v", err)
	}

	cfg.addHeaders(req)
	return req, nil
}

func (cfg *Cfg) validate() error {
	if cfg.URL == "" {
		return errors.New("no endpoint")
	}

	if !hasProtoScheme(cfg.URL) {
		cfg.URL = "https://" + cfg.URL
	}

	if !supported(cfg.Method) {
		return fmt.Errorf("%s is an invalid request method", cfg.Method)
	}
	return nil
}

func (cfg *Cfg) body() (bytes.Buffer, error) {
	var contentType string
	var buf bytes.Buffer

	if cfg.File == "" {
		return buf, nil
	}

	switch ext := path.Ext(cfg.File); ext {
	case jsonExt:
		contentType = "Content-Type:application/json"
	case xmlExt:
		contentType = "Content-Type:application/xml"
	default:
		return buf, fmt.Errorf("%s is an unsupported format", ext)
	}

	cfg.Headers = append(cfg.Headers, contentType)

	b, err := ioutil.ReadFile(cfg.File)
	if err != nil {
		return buf, err
	}
	return *bytes.NewBuffer(b), nil
}

func (cfg *Cfg) addHeaders(r *http.Request) {
	for _, h := range cfg.Headers {
		k := strings.Split(h, ":")[0]
		v := strings.Split(h, ":")[1]
		r.Header.Set(k, v)
		flog.Info("set %s", h)
	}
}

func supported(method string) bool {
	for _, m := range u.Methods {
		if m == method {
			return true
		}
	}
	return false
}

func hasProtoScheme(url string) bool {
	has := func(s string) bool {
		return strings.Contains(url, s)
	}
	return has("http://") || has("https://")
}
