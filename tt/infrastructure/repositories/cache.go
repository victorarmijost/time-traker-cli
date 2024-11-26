package repositories

import (
	"sync"
)

type dbCache struct {
	records map[string]any
	mu      *sync.RWMutex
}

var cache *dbCache
var once sync.Once

func newDBCache() *dbCache {
	once.Do(func() {
		if cache != nil {
			return
		}

		c := &dbCache{
			records: make(map[string]any),
			mu:      &sync.RWMutex{},
		}

		cache = c

	})

	return cache
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

func withCache[T any](c *dbCache, key string, f func() (T, error)) (T, error) {
	var value T
	var err error

	if getFromCache(c, key, &value) {
		return value, nil
	}

	value, err = f()
	if err != nil {
		return value, err
	}

	setInCache(c, key, value)

	return value, nil
}

func withCacheMust[T any](c *dbCache, key string, f func() T) T {
	value, _ := withCache(c, key, func() (T, error) {
		return f(), nil
	})

	return value
}

func resetCache(c *dbCache) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.records = make(map[string]any)
}

func withResetCache(c *dbCache, f func() error) error {
	err := f()
	if err != nil {
		return err
	}

	resetCache(c)

	return nil
}
