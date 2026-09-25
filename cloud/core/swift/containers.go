package swift

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) ListContainers(ctx context.Context) ([]Container, error) {
	url := fmt.Sprintf("%s", strings.TrimRight(c.baseURL, "/"))
	req, _ := http.NewRequest("GET", url, nil)

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var containers []Container
	if err := json.NewDecoder(resp.Body).Decode(&containers); err != nil {
		return nil, err
	}

	return containers, nil
}

func (c *Client) CreateContainer(ctx context.Context, name string) error {
	url := fmt.Sprintf("%s/%s", strings.TrimRight(c.baseURL, "/"), name)
	req, _ := http.NewRequest("PUT", url, nil)

	_, err := c.do(ctx, req)
	return err
}

func (c *Client) DeleteContainer(ctx context.Context, name string) error {
	url := fmt.Sprintf("%s/%s", strings.TrimRight(c.baseURL, "/"), name)
	req, _ := http.NewRequest("DELETE", url, nil)

	_, err := c.do(ctx, req)
	return err
}
