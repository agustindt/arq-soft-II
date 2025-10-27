package services

import (
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/karlseguin/ccache/v3"
)

type DualCache struct {
	local *ccache.Cache[string]
	dist  *memcache.Client
	ttl   time.Duration
}

func NewDualCache(memcachedAddr string, ttlSec int) *DualCache {
	c := ccache.New[string](ccache.Configure[string]().MaxSize(10000))
	mc := memcache.New(memcachedAddr)
	if ttlSec <= 0 {
		ttlSec = 60
	}
	return &DualCache{local: c, dist: mc, ttl: time.Duration(ttlSec) * time.Second}
}

func (d *DualCache) Get(key string) (string, bool) {
	if it := d.local.Get(key); it != nil && !it.Expired() {
		return it.Value(), true
	}
	if v, err := d.dist.Get(key); err == nil {
		d.local.Set(key, string(v.Value), d.ttl)
		return string(v.Value), true
	}
	return "", false
}

func (d *DualCache) Set(key, val string) {
	d.local.Set(key, val, d.ttl)
	_ = d.dist.Set(&memcache.Item{Key: key, Value: []byte(val), Expiration: int32(d.ttl.Seconds())})
}

func (d *DualCache) InvalidatePrefix(prefix string) { d.local.DeletePrefix(prefix) }
