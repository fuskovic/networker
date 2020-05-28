package request

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Run executes the command logic for the request package.
func Run(cfg *Config) error {
	seconds := time.Duration(cfg.TimeOut)
	timeOut := time.Duration(seconds * time.Second)
	client := http.Client{Timeout: timeOut}

	body, err := cfg.buildBody()
	if err != nil {
		return err
	}

	if !cfg.hasProtoScheme() {
		cfg.URL = "https://" + cfg.URL
	}

	if !cfg.validMethod() {
		return fmt.Errorf("%s is an invalid request method", cfg.Method)
	}

	req, err := http.NewRequest(cfg.Method, cfg.URL, body)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
	return nil
}
