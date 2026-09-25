package neutron

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) ListSubnets(ctx context.Context) ([]Subnet, error) {
	url := fmt.Sprintf("%s/v2.0/subnets", strings.TrimRight(c.baseURL, "/"))
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
		Subnets []Subnet `json:"subnets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload.Subnets, nil
}

type CreateSubnetRequest struct {
	Name      string
	NetworkID string
	CIDR      string
	IPVersion int
}

func (c *Client) CreateSubnet(ctx context.Context, reqData CreateSubnetRequest) (*Subnet, error) {
	url := fmt.Sprintf("%s/v2.0/subnets", strings.TrimRight(c.baseURL, "/"))

	body := struct {
		Subnet struct {
			Name      string `json:"name"`
			NetworkID string `json:"network_id"`
			CIDR      string `json:"cidr"`
			IPVersion int    `json:"ip_version"`
		} `json:"subnet"`
	}{}
	body.Subnet.Name = reqData.Name
	body.Subnet.NetworkID = reqData.NetworkID
	body.Subnet.CIDR = reqData.CIDR
	body.Subnet.IPVersion = reqData.IPVersion

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
		Subnet Subnet `json:"subnet"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &payload.Subnet, nil
}
