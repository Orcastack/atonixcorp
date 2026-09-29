package io

import "fmt"

// CacheEngine handles dataset caching for fast access.
type CacheEngine struct {
	cache map[string][]byte
}

func NewCacheEngine() *CacheEngine {
	return &CacheEngine{
		cache: make(map[string][]byte),
	}
}

func (c *CacheEngine) Put(dataset string, data []byte) {
	c.cache[dataset] = data
	fmt.Printf("Cached dataset: %s (%d bytes)\n", dataset, len(data))
}

func (c *CacheEngine) Get(dataset string) ([]byte, bool) {
	data, ok := c.cache[dataset]
	return data, ok
}

func (c *CacheEngine) Evict(dataset string) {
	delete(c.cache, dataset)
	fmt.Printf("Cache evicted: %s\n", dataset)
}

func (c *CacheEngine) Clear() {
	c.cache = make(map[string][]byte)
	fmt.Println("Cache cleared")
}
