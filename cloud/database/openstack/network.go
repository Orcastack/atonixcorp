package openstack

import (
	"encoding/json"
)

type NetworkClient struct {
	os *Client
}

func NewNetworkClient(os *Client) *NetworkClient {
	return &NetworkClient{os}
}

func (n *NetworkClient) CreatePort(networkID, sgID string) (string, error) {
	payload := map[string]interface{}{
		"port": map[string]interface{}{
			"network_id": networkID,
			"security_groups": []string{
				sgID,
			},
		},
	}

	data, err := n.os.do("POST", n.os.AuthURL+"/ports", payload)
	if err != nil {
		return "", err
	}

	var resp struct {
		Port struct {
			ID string `json:"id"`
		} `json:"port"`
	}
	json.Unmarshal(data, &resp)

	return resp.Port.ID, nil
}
