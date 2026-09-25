package object

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SwiftDriver struct {
	Endpoint  string
	Token     string
	Container string
}

func NewSwiftDriver(endpoint, token, container string) *SwiftDriver {
	return &SwiftDriver{Endpoint: endpoint, Token: token, Container: container}
}

func (s *SwiftDriver) Put(key string, data []byte) error {
	url := fmt.Sprintf("%s/%s/%s", s.Endpoint, s.Container, key)
	req, _ := http.NewRequest("PUT", url, bytes.NewReader(data))
	req.Header.Set("X-Auth-Token", s.Token)
	_, err := http.DefaultClient.Do(req)
	return err
}

func (s *SwiftDriver) Get(key string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s", s.Endpoint, s.Container, key)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Auth-Token", s.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func (s *SwiftDriver) Delete(key string) error {
	url := fmt.Sprintf("%s/%s/%s", s.Endpoint, s.Container, key)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("X-Auth-Token", s.Token)
	_, err := http.DefaultClient.Do(req)
	return err
}
