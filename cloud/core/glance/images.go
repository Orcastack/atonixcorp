package glance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) ListImages(ctx context.Context) ([]Image, error) {
	url := fmt.Sprintf("%s/images", strings.TrimRight(c.baseURL, "/"))
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
		Images []Image `json:"images"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload.Images, nil
}

func (c *Client) GetImage(ctx context.Context, id string) (*Image, error) {
	url := fmt.Sprintf("%s/images/%s", strings.TrimRight(c.baseURL, "/"), id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var img Image
	if err := json.NewDecoder(resp.Body).Decode(&img); err != nil {
		return nil, err
	}

	return &img, nil
}

type CreateImageRequest struct {
	Name            string
	DiskFormat      string
	ContainerFormat string
	Visibility      string
}

func (c *Client) CreateImage(ctx context.Context, reqData CreateImageRequest) (*Image, error) {
	url := fmt.Sprintf("%s/images", strings.TrimRight(c.baseURL, "/"))

	body := struct {
		Name            string `json:"name"`
		DiskFormat      string `json:"disk_format"`
		ContainerFormat string `json:"container_format"`
		Visibility      string `json:"visibility"`
	}{
		Name:            reqData.Name,
		DiskFormat:      reqData.DiskFormat,
		ContainerFormat: reqData.ContainerFormat,
		Visibility:      reqData.Visibility,
	}

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

	var payload Image
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

func (c *Client) DeleteImage(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/images/%s", strings.TrimRight(c.baseURL, "/"), id)
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
