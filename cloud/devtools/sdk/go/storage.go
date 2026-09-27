package sdk

import (
	"context"
	"encoding/json"
	"fmt"
)

type Volume struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	SizeGB int    `json:"size_gb"`
	State  string `json:"state"`
}

type CreateVolumeRequest struct {
	Name   string `json:"name"`
	SizeGB int    `json:"size_gb"`
}

func (c *Client) CreateStorage(ctx context.Context, req CreateVolumeRequest) (*Volume, error) {
	body, status, err := c.do(ctx, "POST", "/storage/create", req)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, parseAPIError(status, body)
	}

	var vol Volume
	if err := json.Unmarshal(body, &vol); err != nil {
		return nil, fmt.Errorf("decode volume: %w", err)
	}
	return &vol, nil
}

func (c *Client) DeleteStorage(ctx context.Context, id string) error {
	body, status, err := c.do(ctx, "DELETE", "/storage/delete/"+id, nil)
	if err != nil {
		return err
	}
	if status >= 400 {
		return parseAPIError(status, body)
	}
	return nil
}

func (c *Client) ListStorage(ctx context.Context) ([]Volume, error) {
	body, status, err := c.do(ctx, "GET", "/storage/list", nil)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, parseAPIError(status, body)
	}

	var vols []Volume
	if err := json.Unmarshal(body, &vols); err != nil {
		return nil, fmt.Errorf("decode volume list: %w", err)
	}
	return vols, nil
}
