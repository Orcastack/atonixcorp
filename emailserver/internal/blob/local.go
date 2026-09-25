package blob

import (
	"context"
	"os"
	"path/filepath"
)

type LocalStore struct {
	base string
}

func NewLocalStore(base string) *LocalStore {
	return &LocalStore{base: base}
}

func (l *LocalStore) Put(ctx context.Context, key string, data []byte) error {
	path := filepath.Join(l.base, key)
	os.MkdirAll(filepath.Dir(path), 0755)
	return os.WriteFile(path, data, 0644)
}

func (l *LocalStore) Get(ctx context.Context, key string) ([]byte, error) {
	path := filepath.Join(l.base, key)
	return os.ReadFile(path)
}

func (l *LocalStore) Delete(ctx context.Context, key string) error {
	path := filepath.Join(l.base, key)
	return os.Remove(path)
}
