package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
)

const jsonExt = ".json"

var (
	errUrlUnset                          = errors.New("no url endpoint specified")
	errProtocolUnset                     = errors.New("no protocol specified in url endpoint")
	errInvalidMultiPartFormDataFormat    = errors.New("invalid multipart form upload format")
	errUnevenNumberOfHeaderKeyValuePairs = errors.New("uneven number of key/value pairs")
	errInvalidRequestMethod              = errors.New("invalid request method")
	errUnsupportedFileExtension          = errors.New("unsupported file extension")
	errNoUploadFilesDesignated           = errors.New("no upload files designated")
	errMultiPartFormNameUnset            = errors.New("no multipart form name specified")
)

// Config is used to determine how we build an HTTP request.
type Config struct {
	Headers       []string
	URL           string
	Method        string
	Body          string
	MultiPartForm string
}

// NewNetworkerCraftedHTTPRequest builds a new HTTP request.
func NewNetworkerCraftedHTTPRequest(cfg *Config) (*http.Request, error) {
	// validate URL and method
	switch {
	case cfg.URL == "":
		return nil, errUrlUnset
	case !isSupported(cfg.Method):
		return nil, errInvalidRequestMethod
	default:
		if !strings.Contains(cfg.URL, "http") && !strings.Contains(cfg.URL, "https") {
			return nil, errProtocolUnset
		}
	}

	// normalize method in case the user inputs a lowercase method
	var normalizedMethod string
	for _, char := range cfg.Method {
		normalizedMethod += strings.ToUpper(string(char))
	}

	// initialize request
	req, err := http.NewRequest(normalizedMethod, cfg.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize new request : %w", err)
	}

	// add body if necessary
	if requiresBody(normalizedMethod) {
		if cfg.MultiPartForm == "" {
			if err := cfg.addBody(req); err != nil {
				return nil, fmt.Errorf("failed to add request body: %w", err)
			}
		}
	}

	// provision multi-part form data upload if necessary
	if cfg.MultiPartForm != "" {
		if err := cfg.addForm(req); err != nil {
			return nil, fmt.Errorf("failed to add multi-part form data to request: %w", err)
		}
	}

	// add headers if necessary
	if len(cfg.Headers) != 0 {
		if err := cfg.addHeaders(req); err != nil {
			return nil, fmt.Errorf("failed to add headers: %w", err)
		}
	}
	return req, nil
}

func (cfg *Config) addBody(r *http.Request) error {
	defer func() {
		cfg.Headers = append(cfg.Headers, "Content-Type:application/json")
	}()

	var body io.Reader
	ext := path.Ext(cfg.Body)
	switch {
	case ext != "" && ext != jsonExt:
		return errUnsupportedFileExtension
	case ext == jsonExt:
		b, err := ioutil.ReadFile(cfg.Body)
		if err != nil {
			return fmt.Errorf("failed to read file %q: %w", cfg.Body, err)
		}
		body = bytes.NewBuffer(b)
	default:
		body = strings.NewReader(cfg.Body)
	}
	tempReq, _ := http.NewRequest(http.MethodPost, "temporary", body)
	r.Body = tempReq.Body
	return nil
}

func (cfg *Config) addForm(r *http.Request) error {
	if !strings.Contains(cfg.MultiPartForm, "=") {
		return errInvalidMultiPartFormDataFormat
	}

	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	defer mw.Close()

	args := strings.Split(cfg.MultiPartForm, "=")
	files := strings.Split(args[1], ",")
	for i, f := range files {
		files[i] = strings.TrimSpace(f)
	}
	if len(files) == 0 || (len(files) == 1 && files[0] == "") {
		return errNoUploadFilesDesignated
	}

	formName := args[0]
	if strings.TrimSpace(formName) == "" {
		return errMultiPartFormNameUnset
	}

	for _, filePath := range files {
		f, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open %q: %w", filePath, err)
		}
		defer f.Close()

		ff, err := mw.CreateFormFile(formName, filePath)
		if err != nil {
			return fmt.Errorf("failed to create form file")
		}

		if _, err := io.Copy(ff, f); err != nil {
			return fmt.Errorf("failed to copy %q into form: %w", filePath, err)
		}
	}
	r.Header.Add("Content-Type", mw.FormDataContentType())
	tempReq, _ := http.NewRequest(http.MethodPost, "temporary", &b)
	r.Body = tempReq.Body
	return nil
}

func (cfg *Config) addHeaders(r *http.Request) error {
	for _, h := range cfg.Headers {
		args := strings.Split(h, ":")
		if (len(args) % 2) != 0 {
			return errUnevenNumberOfHeaderKeyValuePairs
		}
		r.Header.Set(args[0], args[1])
	}
	return nil
}

func requiresBody(method string) bool {
	return method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch
}

func isSupported(method string) bool {
	for _, m := range []string{
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
		http.MethodHead,
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	} {
		if strings.ToLower(m) == strings.ToLower(method) {
			return true
		}
	}
	return false
}
