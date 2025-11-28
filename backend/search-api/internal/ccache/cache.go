package ccache

import (
	"sync"
	"time"
)

type Item[T any] struct {
	value      T
	expiration time.Time
}

func (i *Item[T]) Expired() bool {
	return !i.expiration.IsZero() && time.Now().After(i.expiration)
}

func (i *Item[T]) Value() T {
	return i.value
}

type Cache[T any] struct {
	mu      sync.RWMutex
	items   map[string]Item[T]
	maxSize int64
}

type Config[T any] struct {
	maxSize int64
}

func Configure[T any]() *Config[T] {
	return &Config[T]{maxSize: 0}
}

func (c *Config[T]) MaxSize(size int64) *Config[T] {
	c.maxSize = size
	return c
}

func New[T any](cfg *Config[T]) *Cache[T] {
	maxSize := cfg.maxSize
	if maxSize <= 0 {
		maxSize = 10_000
	}
	return &Cache[T]{
		items:   make(map[string]Item[T]),
		maxSize: maxSize,
	}
}

func (c *Cache[T]) Set(key string, value T, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if int64(len(c.items)) >= c.maxSize {
		for k := range c.items {
			delete(c.items, k)
			break
		}
	}

	c.items[key] = Item[T]{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

func (c *Cache[T]) Get(key string) *Item[T] {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		return nil
	}

	if item.Expired() {
		c.Delete(key)
		return nil
	}

	return &item
}

func (c *Cache[T]) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

func (c *Cache[T]) Clear() {
	c.mu.Lock()
	c.items = make(map[string]Item[T])
	c.mu.Unlock()
}
