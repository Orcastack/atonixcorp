package telemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) ListMetrics(ctx context.Context) ([]Metric, error) {
	url := fmt.Sprintf("%s/v1/metric", strings.TrimRight(c.gnocURL, "/"))
	req, _ := http.NewRequest("GET", url, nil)

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var metrics []Metric
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		return nil, err
	}

	return metrics, nil
}
