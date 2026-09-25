package blob

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
)

type GCSStore struct {
	client *storage.Client
	bucket string
}

func NewGCSStore(client *storage.Client, bucket string) *GCSStore {
	return &GCSStore{client: client, bucket: bucket}
}

func (g *GCSStore) Put(ctx context.Context, key string, data []byte) error {
	w := g.client.Bucket(g.bucket).Object(key).NewWriter(ctx)
	_, err := w.Write(data)
	if err != nil {
		return err
	}
	return w.Close()
}

func (g *GCSStore) Get(ctx context.Context, key string) ([]byte, error) {
	r, err := g.client.Bucket(g.bucket).Object(key).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

func (g *GCSStore) Delete(ctx context.Context, key string) error {
	return g.client.Bucket(g.bucket).Object(key).Delete(ctx)
}
