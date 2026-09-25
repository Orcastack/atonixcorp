package swift

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
}

type Config struct {
	Region    string
	Interface string
}

func NewClient(ks *keystone.Client, cfg Config) (*Client, error) {
	ep, err := ks.Endpoint("object-store", cfg.Interface, cfg.Region)
	if err != nil {
		return nil, fmt.Errorf("swift: endpoint lookup failed: %w", err)
	}

	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		ks:         ks,
		baseURL:    ep,
	}, nil
}

func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if err := c.ks.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("X-Auth-Token", c.ks.Token().ID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("swift: http error: %s", resp.Status)
	}

	return resp, nil
}
