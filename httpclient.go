package main

import (
	"io"
	"net/http"
)

type httpClient struct {
	c        http.Client
	apiToken string
}

func (c *httpClient) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err

	}

	return c.Do(req)

}

func (c *httpClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err

	}

	req.Header.Set("Content-Type", contentType)

	return c.Do(req)

}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("PRIVATE-TOKEN", c.apiToken)
	return c.c.Do(req)

}
