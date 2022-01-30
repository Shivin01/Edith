package client

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	httpClient    *http.Client
	getHTTPClient sync.Once
)

// GetHTTPClient returns the http client for this edith to use the default go-client with a Timeout of 10s
func GetHTTPClient() *http.Client {
	getHTTPClient.Do(func() {
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.MaxConnsPerHost = 3

		httpClient = &http.Client{
			Timeout:   time.Second * 10,
			Transport: &botTransport{transport},
		}
	})

	return httpClient
}

// custom http.Transport to set a custom user-agent
type botTransport struct {
	roundTripper http.RoundTripper
}

func (t *botTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	userAgent := fmt.Sprintf("edith/%s", "beta")
	req.Header.Add("User-Agent", userAgent)

	return t.roundTripper.RoundTrip(req)
}
