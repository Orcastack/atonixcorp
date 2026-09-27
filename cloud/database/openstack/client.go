package openstack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type AuthConfig struct {
	AuthURL   string
	Username  string
	Password  string
	Domain    string
	ProjectID string
}

type Client struct {
	Token     string
	ProjectID string
	AuthURL   string
	http      *http.Client
}

func NewClient(cfg AuthConfig) (*Client, error) {
	c := &Client{
		AuthURL:   cfg.AuthURL,
		ProjectID: cfg.ProjectID,
		http:      &http.Client{},
	}

	if err := c.authenticate(cfg); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) authenticate(cfg AuthConfig) error {
	payload := map[string]interface{}{
		"auth": map[string]interface{}{
			"identity": map[string]interface{}{
				"methods": []string{"password"},
				"password": map[string]interface{}{
					"user": map[string]interface{}{
						"name":     cfg.Username,
						"password": cfg.Password,
						"domain": map[string]string{
							"name": cfg.Domain,
						},
					},
				},
			},
			"scope": map[string]interface{}{
				"project": map[string]interface{}{
					"id": cfg.ProjectID,
				},
			},
		},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", c.AuthURL+"/auth/tokens", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("keystone auth failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return fmt.Errorf("keystone auth failed: status %d", resp.StatusCode)
	}

	c.Token = resp.Header.Get("X-Subject-Token")
	return nil
}

func (c *Client) do(method, url string, payload interface{}) ([]byte, error) {
	var body []byte
	if payload != nil {
		body, _ = json.Marshal(payload)
	}

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", c.Token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("openstack error %d: %s", resp.StatusCode, string(data))
	}

	return data, nil
}
