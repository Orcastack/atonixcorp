package openstack

import (
	"encoding/json"
	"fmt"
)

type ComputeClient struct {
	os *Client
}

func NewComputeClient(os *Client) *ComputeClient {
	return &ComputeClient{os}
}

func (c *ComputeClient) CreateServer(name string, cpu int, memoryMB int, networkID string, sgID string) (string, error) {
	payload := map[string]interface{}{
		"server": map[string]interface{}{
			"name":      name,
			"flavorRef": fmt.Sprintf("cpu-%d-mem-%d", cpu, memoryMB),
			"networks": []map[string]string{
				{"uuid": networkID},
			},
			"security_groups": []map[string]string{
				{"name": sgID},
			},
			"imageRef": "db-image", // your DB VM image
		},
	}

	data, err := c.os.do("POST", c.os.AuthURL+"/servers", payload)
	if err != nil {
		return "", err
	}

	var resp struct {
		Server struct {
			ID string `json:"id"`
		} `json:"server"`
	}
	json.Unmarshal(data, &resp)

	return resp.Server.ID, nil
}

func (c *ComputeClient) GetServerIP(serverID string) (string, error) {
	data, err := c.os.do("GET", c.os.AuthURL+"/servers/"+serverID, nil)
	if err != nil {
		return "", err
	}

	var resp struct {
		Server struct {
			Addresses map[string][]struct {
				Addr string `json:"addr"`
			} `json:"addresses"`
		} `json:"server"`
	}
	json.Unmarshal(data, &resp)

	for _, nets := range resp.Server.Addresses {
		for _, addr := range nets {
			return addr.Addr, nil
		}
	}

	return "", fmt.Errorf("no IP found for server %s", serverID)
}
