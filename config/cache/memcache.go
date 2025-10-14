package cache

import (
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// Cache estructura principal para manejar la conexión
type Cache struct {
	Client *memcache.Client
}

// New crea una nueva conexión a Memcached
func New(server string) (*Cache, error) {
	mc := memcache.New(server)
	// Verificamos conexión inicial
	err := mc.Set(&memcache.Item{Key: "ping", Value: []byte("pong"), Expiration: 5})
	if err != nil {
		return nil, err
	}

	log.Println(" Conexión exitosa con Memcached:", server)
	return &Cache{Client: mc}, nil
}

// Set guarda un valor en caché
func (c *Cache) Set(key string, value []byte, expiration time.Duration) error {
	return c.Client.Set(&memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: int32(expiration.Seconds()),
	})
}

// Get obtiene un valor desde caché
func (c *Cache) Get(key string) ([]byte, error) {
	item, err := c.Client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			log.Printf(" Clave no encontrada en caché: %s\n", key)
		}
		return nil, err
	}
	return item.Value, nil
}

// Delete elimina una clave del caché
func (c *Cache) Delete(key string) error {
	return c.Client.Delete(key)
}
