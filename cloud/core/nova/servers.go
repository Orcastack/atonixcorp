package nova

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) ListServers(ctx context.Context) ([]Server, error) {
	url := fmt.Sprintf("%s/servers", strings.TrimRight(c.baseURL, "/"))
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
		Servers []Server `json:"servers"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload.Servers, nil
}

func (c *Client) CreateServer(ctx context.Context, reqData CreateServerRequest) (*Server, error) {
	url := fmt.Sprintf("%s/servers", strings.TrimRight(c.baseURL, "/"))

	body := struct {
		Server struct {
			Name      string `json:"name"`
			FlavorRef string `json:"flavorRef"`
			ImageRef  string `json:"imageRef"`
			Networks  []struct {
				UUID string `json:"uuid"`
			} `json:"networks"`
		} `json:"server"`
	}{}

	body.Server.Name = reqData.Name
	body.Server.FlavorRef = reqData.FlavorRef
	body.Server.ImageRef = reqData.ImageRef
	body.Server.Networks = []struct {
		UUID string `json:"uuid"`
	}{{UUID: reqData.NetworkID}}

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
		Server Server `json:"server"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &payload.Server, nil
}

func (c *Client) DeleteServer(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/servers/%s", strings.TrimRight(c.baseURL, "/"), id)
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
