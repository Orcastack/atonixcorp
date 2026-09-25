package heat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) ListStacks(ctx context.Context) ([]Stack, error) {
	url := fmt.Sprintf("%s/stacks", strings.TrimRight(c.baseURL, "/"))
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
		Stacks []Stack `json:"stacks"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload.Stacks, nil
}

func (c *Client) CreateStack(ctx context.Context, reqData CreateStackRequest) (*Stack, error) {
	url := fmt.Sprintf("%s/stacks", strings.TrimRight(c.baseURL, "/"))

	body := struct {
		StackName  string            `json:"stack_name"`
		Template   string            `json:"template"`
		Parameters map[string]string `json:"parameters,omitempty"`
	}{
		StackName:  reqData.Name,
		Template:   reqData.Template,
		Parameters: reqData.Parameters,
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

	var payload Stack
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

func (c *Client) DeleteStack(ctx context.Context, name, id string) error {
	// Heat identifies stacks by name + ID
	url := fmt.Sprintf("%s/stacks/%s/%s", strings.TrimRight(c.baseURL, "/"), name, id)
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
