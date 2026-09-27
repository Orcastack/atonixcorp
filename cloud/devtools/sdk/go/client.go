package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

type Client struct {
	endpoint   string
	token      string
	httpClient *http.Client
}

type ClientOptions struct {
	Timeout      time.Duration
	MaxIdleConns int
	IdleConnTTL  time.Duration
	Endpoint     string
	Token        string
}

func NewClient(opts ClientOptions) *Client {
	if opts.Timeout == 0 {
		opts.Timeout = 30 * time.Second
	}
	if opts.MaxIdleConns == 0 {
		opts.MaxIdleConns = 100
	}
	if opts.IdleConnTTL == 0 {
		opts.IdleConnTTL = 90 * time.Second
	}

	transport := &http.Transport{
		MaxIdleConns:        opts.MaxIdleConns,
		IdleConnTimeout:     opts.IdleConnTTL,
		TLSHandshakeTimeout: 10 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	return &Client{
		endpoint: opts.Endpoint,
		token:    opts.Token,
		httpClient: &http.Client{
			Timeout:   opts.Timeout,
			Transport: transport,
		},
	}
}

func (c *Client) do(ctx context.Context, method, path string, body any) ([]byte, int, error) {
	var buf io.Reader

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal request: %w", err)
		}
		buf = bytes.NewBuffer(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.endpoint+path, buf)
	if err != nil {
		return nil, 0, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("http error: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("read response: %w", err)
	}

	return data, res.StatusCode, nil
}
