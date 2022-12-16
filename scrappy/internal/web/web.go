// Package web handles http requests and data scraping.
package web

import (
	"net/http"
	"time"
)

const defaultTimeout time.Duration = 10 * time.Second

// NewClient returns a new HTTP client with the given timeout.
//
// If timeout is 0, it will use the `defaultTimeout`
func NewClient(timeout time.Duration) *http.Client {
	if timeout == 0 {
		timeout = defaultTimeout
	}
	return &http.Client{Timeout: timeout}
}

// CheckURL send an http HEAD request to the url to check if it is reachable.
func CheckURL(url string) (status int, err error) {
	response, err := NewClient(defaultTimeout).Head(url)
	if err != nil {
		return 0, err
	}

	return response.StatusCode, nil
}
