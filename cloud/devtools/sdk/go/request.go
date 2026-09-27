package sdk

import (
	"encoding/json"
	"fmt"
)

type APIError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api error (%d): %s - %s", e.StatusCode, e.Code, e.Message)
}

func parseAPIError(status int, body []byte) error {
	var apiErr APIError
	if err := json.Unmarshal(body, &apiErr); err != nil {
		return fmt.Errorf("status %d: %s", status, string(body))
	}
	apiErr.StatusCode = status
	return &apiErr
}
