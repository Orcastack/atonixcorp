package sdk

import (
	"context"
	"encoding/json"
	"fmt"
)

type Event struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type Metric struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func (c *Client) ListEvents(ctx context.Context) ([]Event, error) {
	body, status, err := c.do(ctx, "GET", "/orchestration/events", nil)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, parseAPIError(status, body)
	}

	var events []Event
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, fmt.Errorf("decode events: %w", err)
	}
	return events, nil
}

func (c *Client) ListMetrics(ctx context.Context) ([]Metric, error) {
	body, status, err := c.do(ctx, "GET", "/orchestration/metrics", nil)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, parseAPIError(status, body)
	}

	var metrics []Metric
	if err := json.Unmarshal(body, &metrics); err != nil {
		return nil, fmt.Errorf("decode metrics: %w", err)
	}
	return metrics, nil
}
