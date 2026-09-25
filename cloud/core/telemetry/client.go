package telemetry

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
	ceilURL    string
	gnocURL    string
	pankoURL   string
}

type Config struct {
	Region    string
	Interface string
}

func NewClient(ks *keystone.Client, cfg Config) (*Client, error) {
	ceil, err := ks.Endpoint("metering", cfg.Interface, cfg.Region)
	if err != nil {
		return nil, fmt.Errorf("telemetry: ceilometer endpoint failed: %w", err)
	}

	gnoc, err := ks.Endpoint("metric", cfg.Interface, cfg.Region)
	if err != nil {
		return nil, fmt.Errorf("telemetry: gnocchi endpoint failed: %w", err)
	}

	panko, err := ks.Endpoint("event", cfg.Interface, cfg.Region)
	if err != nil {
		return nil, fmt.Errorf("telemetry: panko endpoint failed: %w", err)
	}

	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		ks:         ks,
		ceilURL:    ceil,
		gnocURL:    gnoc,
		pankoURL:   panko,
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
		return nil, fmt.Errorf("telemetry: http error: %s", resp.Status)
	}

	return resp, nil
}
