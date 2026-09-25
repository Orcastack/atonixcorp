package neutron

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) ListPorts(ctx context.Context) ([]Port, error) {
	url := fmt.Sprintf("%s/v2.0/ports", strings.TrimRight(c.baseURL, "/"))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var payload struct {
		Ports []Port `json:"ports"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload.Ports, nil
}
