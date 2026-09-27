package openstack

import (
	"encoding/json"
	"fmt"
)

type VolumeClient struct {
	os *Client
}

func NewVolumeClient(os *Client) *VolumeClient {
	return &VolumeClient{os}
}

func (v *VolumeClient) CreateVolume(name string, sizeGB int) (string, error) {
	payload := map[string]interface{}{
		"volume": map[string]interface{}{
			"name": name,
			"size": sizeGB,
		},
	}

	data, err := v.os.do("POST", v.os.AuthURL+"/volumes", payload)
	if err != nil {
		return "", err
	}

	var resp struct {
		Volume struct {
			ID string `json:"id"`
		} `json:"volume"`
	}
	json.Unmarshal(data, &resp)

	return resp.Volume.ID, nil
}

func (v *VolumeClient) AttachVolume(serverID, volumeID string) error {
	payload := map[string]interface{}{
		"volumeAttachment": map[string]interface{}{
			"volumeId": volumeID,
		},
	}

	_, err := v.os.do("POST", fmt.Sprintf("%s/servers/%s/os-volume_attachments", v.os.AuthURL, serverID), payload)
	return err
}
