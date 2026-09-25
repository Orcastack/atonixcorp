package object

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type S3Driver struct {
	Endpoint string
	Bucket   string
}

func NewS3Driver(endpoint, bucket string) *S3Driver {
	return &S3Driver{Endpoint: endpoint, Bucket: bucket}
}

func (s *S3Driver) Put(key string, data []byte) error {
	url := fmt.Sprintf("%s/%s/%s", s.Endpoint, s.Bucket, key)
	_, err := http.Put(url, "application/octet-stream", bytes.NewReader(data))
	return err
}

func (s *S3Driver) Get(key string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s", s.Endpoint, s.Bucket, key)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func (s *S3Driver) Delete(key string) error {
	url := fmt.Sprintf("%s/%s/%s", s.Endpoint, s.Bucket, key)
	req, _ := http.NewRequest("DELETE", url, nil)
	_, err := http.DefaultClient.Do(req)
	return err
}
