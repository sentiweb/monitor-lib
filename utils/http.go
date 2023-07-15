package utils

import (
	"crypto/tls"
	"net/http"
	"time"
)

const (
	// HttpHeaderContentType  HTTP header for Content Type
	HttpHeaderContentType = "Content-Type"

	// MimeJson is the MIME type for JSON
	MimeJson = "application/json"
)

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HTTPClientParams struct {
	insecure bool
}

// HTTPClientFactory creates HTTPClient
type HTTPClientFactory interface {
	NewClient(timeout time.Duration, config *HTTPClientParams) HTTPClient
}

type DefaultHTTPFactory struct{}

func (f *DefaultHTTPFactory) NewClient(timeout time.Duration, config *HTTPClientParams) HTTPClient {
	if config == nil {
		config = &HTTPClientParams{}
	}
	var client *http.Client
	if config.insecure {
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client = &http.Client{Transport: customTransport, Timeout: timeout}
	} else {
		client = &http.Client{Timeout: timeout}
	}
	return client
}
