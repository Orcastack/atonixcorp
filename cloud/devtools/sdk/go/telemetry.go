package sdk

import (
	"context"
	"encoding/json"
	"fmt"
)

type TelemetrySnapshot struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
	Timestamp   int64   `json:"timestamp"`
}

func (c *Client) Telemetry(ctx context.Context) (*TelemetrySnapshot, error) {
	body, status, err := c.do(ctx, "GET", "/telemetry", nil)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, parseAPIError(status, body)
	}

	var snap TelemetrySnapshot
	if err := json.Unmarshal(body, &snap); err != nil {
		return nil, fmt.Errorf("decode telemetry: %w", err)
	}
	return &snap, nil
}
