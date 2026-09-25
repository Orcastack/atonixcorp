package object

import (
	"bytes"
	"context"
	"io/ioutil"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioDriver struct {
	client *minio.Client
	bucket string
}

func NewMinioDriver(endpoint, accessKey, secretKey, bucket string, secure bool) (*MinioDriver, error) {
	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, err
	}

	return &MinioDriver{client: cli, bucket: bucket}, nil
}

func (m *MinioDriver) Put(key string, data []byte) error {
	_, err := m.client.PutObject(context.Background(), m.bucket, key,
		bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	return err
}

func (m *MinioDriver) Get(key string) ([]byte, error) {
	obj, err := m.client.GetObject(context.Background(), m.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(obj)
}

func (m *MinioDriver) Delete(key string) error {
	return m.client.RemoveObject(context.Background(), m.bucket, key, minio.RemoveObjectOptions{})
}
