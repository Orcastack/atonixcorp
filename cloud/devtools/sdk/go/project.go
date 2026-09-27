package sdk

import (
	"context"
	"encoding/json"
	"fmt"
)

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) ListProjects(ctx context.Context) ([]Project, error) {
	body, status, err := c.do(ctx, "GET", "/projects/list", nil)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, parseAPIError(status, body)
	}

	var projects []Project
	if err := json.Unmarshal(body, &projects); err != nil {
		return nil, fmt.Errorf("decode projects: %w", err)
	}
	return projects, nil
}
