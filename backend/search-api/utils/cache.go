package utils

import (
	"log"
	"sync"
	"time"
)

// SimpleCache es una implementaciﾃｳn en memoria simple del cachﾃｩ
type Cache struct {
	data sync.Map
}

type cacheItem struct {
	value      []byte
	expiration time.Time
}

// New crea una nueva instancia de cachﾃｩ
func NewCache(server string) (*Cache, error) {
	log.Println("Cache en memoria inicializado (server:", server, ")")
	return &Cache{}, nil
}

// Set guarda un valor en cachﾃｩ
func (c *Cache) Set(key string, value []byte, expiration time.Duration) error {
	item := cacheItem{
		value:      value,
		expiration: time.Now().Add(expiration),
	}
	c.data.Store(key, item)
	return nil
}

// Get obtiene un valor desde cachﾃｩ
func (c *Cache) Get(key string) ([]byte, error) {
	val, ok := c.data.Load(key)
	if !ok {
		log.Printf("Clave no encontrada en cachﾃｩ: %s\n", key)
		return nil, ErrCacheMiss
	}

	item := val.(cacheItem)

	// Verificar expiraciﾃｳn
	if time.Now().After(item.expiration) {
		c.data.Delete(key)
		log.Printf("Clave expirada en cachﾃｩ: %s\n", key)
		return nil, ErrCacheMiss
	}

	return item.value, nil
}

// Delete elimina una clave del cachﾃｩ
func (c *Cache) Delete(key string) error {
	c.data.Delete(key)
	return nil
}

// FlushAll borra todo el cachﾃｩ
func (c *Cache) FlushAll() error {
	c.data.Range(func(key, value interface{}) bool {
		c.data.Delete(key)
		return true
	})
	return nil
}

// ErrCacheMiss indica que la clave no fue encontrada
var ErrCacheMiss = &cacheError{message: "cache miss"}

type cacheError struct {
	message string
}

func (e *cacheError) Error() string {
	return e.message
}
