package cinder

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"atonixcorp/cloud/core/keystone"
)

type Client struct {
	httpClient *http.Client
	ks         *keystone.Client
	baseURL    string
	region     string
	iface      string
}

type Config struct {
	Region    string
	Interface string // "public", "internal", "admin"
}

func NewClient(ks *keystone.Client, cfg Config) (*Client, error) {
	ep, err := ks.Endpoint("volumev3", cfg.Interface, cfg.Region)
	if err != nil {
		return nil, fmt.Errorf("cinder: get endpoint: %w", err)
	}

	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		ks:         ks,
		baseURL:    ep,
		region:     cfg.Region,
		iface:      cfg.Interface,
	}, nil
}

func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if err := c.ks.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	if tok := c.ks.Token(); tok != nil {
		req.Header.Set("X-Auth-Token", tok.ID)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("cinder: http error: %s", resp.Status)
	}

	return resp, nil
}
