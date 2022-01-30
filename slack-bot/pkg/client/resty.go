package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

type APIClient struct {
	HTTPClient *resty.Client
}

func NewAPIClient(baseURL string, maxConnsPerHost int, timeout time.Duration) *APIClient {
	apiClient := &APIClient{
		HTTPClient: resty.New(),
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxConnsPerHost = maxConnsPerHost
	apiClient.HTTPClient.
		SetBaseURL(baseURL).
		SetTimeout(timeout * time.Second)

	for key, value := range headers {
		apiClient.HTTPClient.Header.Set(key, value)
	}

	return apiClient
}

func (c *APIClient) SetAuthToken(token string) {
	c.HTTPClient.SetAuthToken(token)
}

func (c *APIClient) post(ctx context.Context, url string, body interface{}, res interface{}, token string, headers map[string]string) (*resty.Response, error) {
	request := c.HTTPClient.R()
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	if token != "" {
		request.SetAuthToken(token)
	}

	if res != nil {
		request.SetResult(res)
	}

	return request.
		SetContext(ctx).
		SetBody(body).
		Post(url)
}

func (c *APIClient) delete(ctx context.Context, url string, token string, headers map[string]string) (*resty.Response, error) {
	request := c.HTTPClient.R()
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	if token != "" {
		request.SetAuthToken(token)
	}

	return request.
		SetContext(ctx).
		Delete(url)
}

func (c *APIClient) Post(ctx context.Context, url string, body interface{}, res interface{}, headers map[string]string) (*resty.Response, error) {
	return c.handle(c.post(ctx, url, body, res, "", headers))
}

func (c *APIClient) PostWithToken(ctx context.Context, url string, body interface{}, res interface{}, token string, headers map[string]string) (*resty.Response, error) {
	return c.handle(c.post(ctx, url, body, res, token, headers))
}

func (c *APIClient) get(ctx context.Context, url string, token string, res interface{}, headers map[string]string) (*resty.Response, error) {
	request := c.HTTPClient.R()
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	if token != "" {
		request.SetAuthToken(token)
	}

	if res != nil {
		request.SetResult(res)
	}

	return request.SetContext(ctx).Get(url)
}

func (c *APIClient) patch(ctx context.Context, url string, token string, body interface{}, res interface{}, headers map[string]string) (*resty.Response, error) {
	request := c.HTTPClient.R()
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	if token != "" {
		request.SetAuthToken(token)
	}

	if res != nil {
		request.SetResult(res)
	}

	return request.SetContext(ctx).SetBody(body).Patch(url)
}

func (c *APIClient) Get(ctx context.Context, url string, res interface{}, headers map[string]string) (*resty.Response, error) {
	return c.handle(c.get(ctx, url, "", res, headers))
}

func (c *APIClient) GetWithToken(ctx context.Context, url, token string, res interface{}, headers map[string]string) (*resty.Response, error) {
	return c.handle(c.get(ctx, url, token, res, headers))
}

func (c *APIClient) Delete(ctx context.Context, url, token string, headers map[string]string) (*resty.Response, error) {
	return c.handle(c.delete(ctx, url, token, headers))
}

func (c *APIClient) Patch(ctx context.Context, url, token string, body interface{}, res interface{}, headers map[string]string) (*resty.Response, error) {
	return c.handle(c.patch(ctx, url, token, body, res, headers))
}

func (c *APIClient) handle(res *resty.Response, err error) (*resty.Response, error) {
	if err != nil {
		return nil, err
	}
	if res.StatusCode() >= 300 {
		return nil, errors.New(fmt.Sprintf("Got %d status code", res.StatusCode()))
	}
	return res, nil
}
