package sdk

import (
	"context"
	"encoding/json"
	"fmt"
)

type Image struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	OS      string `json:"os"`
	Version string `json:"version"`
}

func (c *Client) ListImages(ctx context.Context) ([]Image, error) {
	body, status, err := c.do(ctx, "GET", "/images/list", nil)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, parseAPIError(status, body)
	}

	var imgs []Image
	if err := json.Unmarshal(body, &imgs); err != nil {
		return nil, fmt.Errorf("decode image list: %w", err)
	}
	return imgs, nil
}
