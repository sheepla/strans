package api

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultInstance = "lingva.ml"
	timeout         = 10 * time.Second
)

func httpGet(req *http.Request) (io.ReadCloser, error) {
	//nolint:exhaustivestruct,exhaustruct
	cl := &http.Client{
		Timeout: timeout,
	}

	resp, err := cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		return nil, fmt.Errorf("%w: %s", ErrHTTP, resp.Status)
	}

	return resp.Body, nil
}
