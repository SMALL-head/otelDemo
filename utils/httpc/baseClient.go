package httpc

import (
	"bytes"
	"net/http"
	"net/url"
	"time"
)

type BaseClient struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func NewClient(base string) (*BaseClient, error) {
	parsedURL, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	return &BaseClient{
		baseURL: parsedURL,
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}, nil
}

func (c *BaseClient) DoGet(path string) (*http.Response, error) {
	fullUrl := c.baseURL.ResolveReference(&url.URL{Path: path})
	return c.httpClient.Get(fullUrl.String())
}

func (c *BaseClient) DoPost(path string, body []byte, contentType string) (*http.Response, error) {
	fullUrl := c.baseURL.ResolveReference(&url.URL{Path: path})
	return c.httpClient.Post(fullUrl.String(), contentType, bytes.NewBuffer(body))
}
