package swift

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) UploadObject(ctx context.Context, container, name string, data io.Reader) error {
	url := fmt.Sprintf("%s/%s/%s", strings.TrimRight(c.baseURL, "/"), container, name)
	req, _ := http.NewRequest("PUT", url, data)

	_, err := c.do(ctx, req)
	return err
}

func (c *Client) DownloadObject(ctx context.Context, container, name string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s/%s/%s", strings.TrimRight(c.baseURL, "/"), container, name)
	req, _ := http.NewRequest("GET", url, nil)

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (c *Client) DeleteObject(ctx context.Context, container, name string) error {
	url := fmt.Sprintf("%s/%s/%s", strings.TrimRight(c.baseURL, "/"), container, name)
	req, _ := http.NewRequest("DELETE", url, nil)

	_, err := c.do(ctx, req)
	return err
}
