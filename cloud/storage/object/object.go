package object

import (
	"fmt"
	"sync"
)

type LocalObjectStore struct {
	store map[string][]byte
	mu    sync.RWMutex
}

func NewLocalObjectStore() *LocalObjectStore {
	return &LocalObjectStore{
		store: make(map[string][]byte),
	}
}

func (o *LocalObjectStore) Put(key string, data []byte) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.store[key] = data
	fmt.Println("LocalObjectStore: stored", key)
	return nil
}

func (o *LocalObjectStore) Get(key string) ([]byte, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	data, ok := o.store[key]
	if !ok {
		return nil, fmt.Errorf("LocalObjectStore: key not found: %s", key)
	}

	fmt.Println("LocalObjectStore: get", key)
	return data, nil
}

func (o *LocalObjectStore) Delete(key string) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	delete(o.store, key)
	fmt.Println("LocalObjectStore: delete", key)
	return nil
}
