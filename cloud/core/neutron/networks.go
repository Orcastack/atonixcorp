package neutron

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) ListNetworks(ctx context.Context) ([]Network, error) {
	url := fmt.Sprintf("%s/v2.0/networks", strings.TrimRight(c.baseURL, "/"))
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
		Networks []Network `json:"networks"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload.Networks, nil
}

type CreateNetworkRequest struct {
	Name         string
	AdminStateUp bool
}

func (c *Client) CreateNetwork(ctx context.Context, reqData CreateNetworkRequest) (*Network, error) {
	url := fmt.Sprintf("%s/v2.0/networks", strings.TrimRight(c.baseURL, "/"))

	body := struct {
		Network struct {
			Name         string `json:"name"`
			AdminStateUp bool   `json:"admin_state_up"`
		} `json:"network"`
	}{}
	body.Network.Name = reqData.Name
	body.Network.AdminStateUp = reqData.AdminStateUp

	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(buf)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var payload struct {
		Network Network `json:"network"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &payload.Network, nil
}

func (c *Client) DeleteNetwork(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/v2.0/networks/%s", strings.TrimRight(c.baseURL, "/"), id)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}
