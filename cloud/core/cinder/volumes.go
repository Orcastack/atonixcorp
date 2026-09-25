package cinder

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) ListVolumes(ctx context.Context) ([]Volume, error) {
	url := fmt.Sprintf("%s/volumes", strings.TrimRight(c.baseURL, "/"))
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
		Volumes []Volume `json:"volumes"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload.Volumes, nil
}

type CreateVolumeRequest struct {
	Name        string
	Size        int
	Description string
}

func (c *Client) CreateVolume(ctx context.Context, reqData CreateVolumeRequest) (*Volume, error) {
	url := fmt.Sprintf("%s/volumes", strings.TrimRight(c.baseURL, "/"))

	body := struct {
		Volume struct {
			Name        string `json:"name"`
			Size        int    `json:"size"`
			Description string `json:"description,omitempty"`
		} `json:"volume"`
	}{}
	body.Volume.Name = reqData.Name
	body.Volume.Size = reqData.Size
	body.Volume.Description = reqData.Description

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
		Volume Volume `json:"volume"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &payload.Volume, nil
}

func (c *Client) DeleteVolume(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/volumes/%s", strings.TrimRight(c.baseURL, "/"), id)
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
