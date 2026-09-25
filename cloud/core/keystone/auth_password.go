package keystone

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type authRequest struct {
	Auth struct {
		Identity struct {
			Methods  []string        `json:"methods"`
			Password passwordPayload `json:"password"`
		} `json:"identity"`
		Scope scopePayload `json:"scope"`
	} `json:"auth"`
}

type passwordPayload struct {
	User struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Domain   struct {
			ID string `json:"id"`
		} `json:"domain"`
	} `json:"user"`
}

type scopePayload struct {
	Project struct {
		ID string `json:"id"`
	} `json:"project"`
}

func (c *Client) AuthenticatePassword(ctx context.Context) error {
	body := authRequest{}
	body.Auth.Identity.Methods = []string{"password"}
	body.Auth.Identity.Password.User.Name = c.creds.Username
	body.Auth.Identity.Password.User.Password = c.creds.Password
	body.Auth.Identity.Password.User.Domain.ID = c.creds.DomainID
	body.Auth.Scope.Project.ID = c.creds.ProjectID

	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/auth/tokens", c.creds.AuthURL)
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("keystone: auth failed: %s", resp.Status)
	}

	tokenID := resp.Header.Get("X-Subject-Token")

	var payload struct {
		Token struct {
			Project struct {
				ID string `json:"id"`
			} `json:"project"`
			User struct {
				ID string `json:"id"`
			} `json:"user"`
			ExpiresAt string    `json:"expires_at"`
			Roles     []Role    `json:"roles"`
			Catalog   []Service `json:"catalog"`
		} `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return err
	}

	c.token = &Token{
		ID:             tokenID,
		ProjectID:      payload.Token.Project.ID,
		UserID:         payload.Token.User.ID,
		ExpiresAt:      payload.Token.ExpiresAt,
		Roles:          payload.Token.Roles,
		ServiceCatalog: payload.Token.Catalog,
	}

	return nil
}
