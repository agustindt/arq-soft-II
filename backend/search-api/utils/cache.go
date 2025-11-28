package utils

import (
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/karlseguin/ccache/v3"
)

// Cache implementa un caché en dos niveles: memoria local (CCache) y Memcached remoto.
type Cache struct {
	local  *ccache.Cache[string]
	remote *memcache.Client
	ttl    time.Duration
}

// NewCache crea una instancia del caché con TTL configurable.
func NewCache(server string, ttl time.Duration) (*Cache, error) {
	if ttl <= 0 {
		ttl = 30 * time.Second
	}

	cache := &Cache{
		local: ccache.New(ccache.Configure[string]().MaxSize(10_000)),
		ttl:   ttl,
	}

	if server != "" {
		cache.remote = memcache.New(server)
		// Sondeo ligero para validar conexión
		if err := cache.remote.Set(&memcache.Item{Key: "cache_ping", Value: []byte("ok"), Expiration: int32(ttl.Seconds())}); err != nil {
			log.Printf("⚠️ Memcached no respondió al ping inicial: %v", err)
		} else {
			log.Printf("✅ Memcached conectado (%s) con TTL %s", server, ttl)
		}
	} else {
		log.Printf("⚠️ Memcached server vacío, solo se usará caché local en memoria")
	}

	return cache, nil
}

// Set guarda un valor en ambos niveles de caché.
func (c *Cache) Set(key string, value []byte, expiration time.Duration) error {
	if expiration <= 0 {
		expiration = c.ttl
	}

	c.local.Set(key, string(value), expiration)

	if c.remote != nil {
		if err := c.remote.Set(&memcache.Item{Key: key, Value: value, Expiration: int32(expiration.Seconds())}); err != nil {
			log.Printf("⚠️ No se pudo guardar en Memcached: %v", err)
		}
	}

	return nil
}

// Get obtiene un valor desde CCache y, si no existe, lo busca en Memcached.
func (c *Cache) Get(key string) ([]byte, error) {
	if item := c.local.Get(key); item != nil && !item.Expired() {
		return []byte(item.Value()), nil
	}

	if c.remote != nil {
		mcItem, err := c.remote.Get(key)
		if err == nil {
			c.local.Set(key, string(mcItem.Value), c.ttl)
			return mcItem.Value, nil
		}

		if err != memcache.ErrCacheMiss {
			log.Printf("⚠️ Error leyendo Memcached: %v", err)
		}
	}

	return nil, ErrCacheMiss
}

// Delete elimina una clave del caché.
func (c *Cache) Delete(key string) error {
	c.local.Delete(key)
	if c.remote != nil {
		if err := c.remote.Delete(key); err != nil && err != memcache.ErrCacheMiss {
			log.Printf("⚠️ No se pudo borrar clave en Memcached: %v", err)
		}
	}
	return nil
}

// FlushAll limpia los cachés locales y remotos.
func (c *Cache) FlushAll() error {
	c.local.Clear()
	if c.remote != nil {
		if err := c.remote.FlushAll(); err != nil {
			log.Printf("⚠️ No se pudo limpiar Memcached: %v", err)
		}
	}
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
