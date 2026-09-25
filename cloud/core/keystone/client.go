package keystone

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	creds      Credentials
	token      *Token
}

func NewClient(creds Credentials) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		creds: creds,
	}
}

func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	if c.token != nil {
		req.Header.Set("X-Auth-Token", c.token.ID)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return nil, ErrUnauthorized
	case http.StatusForbidden:
		return nil, ErrForbidden
	case http.StatusNotFound:
		return nil, ErrNotFound
	case http.StatusBadRequest:
		return nil, ErrBadRequest
	}

	if resp.StatusCode >= 500 {
		return nil, ErrServerError
	}

	return resp, nil
}

func (c *Client) Token() *Token {
	return c.token
}

func (c *Client) ServiceCatalog() []Service {
	if c.token == nil {
		return nil
	}
	return c.token.ServiceCatalog
}

func (c *Client) ensureToken(ctx context.Context) error {
	if c.token == nil {
		return c.AuthenticatePassword(ctx)
	}
	// simple expiry check; you can parse ExpiresAt properly later
	return nil
}

func decodeJSON(resp *http.Response, v any) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}
