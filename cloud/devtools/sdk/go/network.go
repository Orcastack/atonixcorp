package sdk

import (
	"context"
	"encoding/json"
	"fmt"
)

type Network struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CIDR      string `json:"cidr"`
	ProjectID string `json:"project_id"`
	State     string `json:"state"`
}

type CreateNetworkRequest struct {
	Name      string `json:"name"`
	CIDR      string `json:"cidr"`
	ProjectID string `json:"project_id,omitempty"`
}

type DeleteNetworkRequest struct {
	ID string `json:"id"`
}

// ------------------------------------------------------------
// CREATE NETWORK
// ------------------------------------------------------------
func (c *Client) CreateNetwork(ctx context.Context, req CreateNetworkRequest) (*Network, error) {
	body, status, err := c.do(ctx, "POST", "/network/create", req)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, parseAPIError(status, body)
	}

	var net Network
	if err := json.Unmarshal(body, &net); err != nil {
		return nil, fmt.Errorf("decode network: %w", err)
	}

	return &net, nil
}

// ------------------------------------------------------------
// DELETE NETWORK
// ------------------------------------------------------------
func (c *Client) DeleteNetwork(ctx context.Context, id string) error {
	body, status, err := c.do(ctx, "DELETE", "/network/delete/"+id, nil)
	if err != nil {
		return err
	}
	if status >= 400 {
		return parseAPIError(status, body)
	}
	return nil
}

// ------------------------------------------------------------
// LIST NETWORKS
// ------------------------------------------------------------
func (c *Client) ListNetwork(ctx context.Context) ([]Network, error) {
	body, status, err := c.do(ctx, "GET", "/network/list", nil)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, parseAPIError(status, body)
	}

	var nets []Network
	if err := json.Unmarshal(body, &nets); err != nil {
		return nil, fmt.Errorf("decode network list: %w", err)
	}

	return nets, nil
}
