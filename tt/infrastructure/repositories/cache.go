package repositories

import (
	"sync"
)

type dbCache struct {
	records map[string]any
	mu      *sync.RWMutex
}

func newDBCache() *dbCache {
	return &dbCache{
		records: make(map[string]any),
		mu:      &sync.RWMutex{},
	}
}

func getFromCache[T any](c *dbCache, key string, value *T) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.records[key]
	if !ok {
		return false
	}

	iv, ok := v.(T)
	if !ok {
		return false
	}

	*value = iv

	return true
}

func setInCache[T any](c *dbCache, key string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.records[key] = value
}

func resetCache(c *dbCache) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.records = make(map[string]any)
}
