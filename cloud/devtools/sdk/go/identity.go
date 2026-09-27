package sdk

import (
	"context"
	"encoding/json"
	"fmt"
)

type IdentityInfo struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Projects []string `json:"projects"`
}

func (c *Client) IdentityInfo(ctx context.Context) (*IdentityInfo, error) {
	body, status, err := c.do(ctx, "GET", "/identity/info", nil)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, parseAPIError(status, body)
	}

	var info IdentityInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("decode identity info: %w", err)
	}
	return &info, nil
}
