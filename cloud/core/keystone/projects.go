package keystone

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) ListProjects(ctx context.Context) ([]Project, error) {
	if err := c.ensureToken(ctx); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/projects", c.creds.AuthURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}

	var payload struct {
		Projects []Project `json:"projects"`
	}
	if err := decodeJSON(resp, &payload); err != nil {
		return nil, err
	}

	return payload.Projects, nil
}
