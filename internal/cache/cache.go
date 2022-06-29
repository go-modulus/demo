package cache

import (
	"github.com/Code-Hex/go-generics-cache/policy/lru"
	"time"
)

type Cache[KeyT comparable, ValueT any] struct {
	cache        *lru.Cache[KeyT, ValueT]
	cacheEnabled bool
}

type Config struct {
	MaxCachedItems int
	CacheEnabled   bool
	// time after which entry can be evicted
	LifeTime time.Duration
}

func NewCache[KeyT comparable, ValueT any](config *Config) *Cache[KeyT, ValueT] {
	if config.CacheEnabled {
		cache := lru.NewCache[KeyT, ValueT](lru.WithCapacity(config.MaxCachedItems))

		return &Cache[KeyT, ValueT]{
			cache:        cache,
			cacheEnabled: true,
		}
	}
	return &Cache[KeyT, ValueT]{}
}

func (c *Cache[KeyT, ValueT]) Get(key KeyT) (ValueT, bool) {
	if c.cacheEnabled {
		if v, ok := c.cache.Get(key); ok {
			return v, true
		} else {
			var v ValueT
			return v, false
		}
	}
	var v ValueT
	return v, false
}

func (c *Cache[KeyT, ValueT]) Set(key KeyT, value ValueT) error {
	if c.cacheEnabled {
		c.cache.Set(key, value)
	}
	return nil
}

func (c *Cache[KeyT, ValueT]) Del(key KeyT) error {
	if c.cacheEnabled {
		c.cache.Delete(key)
	}
	return nil
}
