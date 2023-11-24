// Package httpx provides helper functions for making HTTP requests.
package httpx

import (
	"net/http"
	"time"
)

var client = Client(&http.Client{Timeout: 5 * time.Second})

// Client is something that can send HTTP requests.
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

// Do sends an HTTP request and returns an HTTP response.
func Do(req *http.Request) (*http.Response, error) {
	return client.Do(req)
}
