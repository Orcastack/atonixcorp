package telemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) ListEvents(ctx context.Context) ([]Event, error) {
	url := fmt.Sprintf("%s/v2/events", strings.TrimRight(c.pankoURL, "/"))
	req, _ := http.NewRequest("GET", url, nil)

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var payload struct {
		Events []Event `json:"events"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload.Events, nil
}
