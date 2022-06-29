package cache

type LoaderCache[KeyT comparable, ValueT any] struct {
	baseCache *Cache[KeyT, ValueT]
}

func NewLoaderCache[KeyT comparable, ValueT any](baseCache *Cache[KeyT, ValueT]) *LoaderCache[KeyT, ValueT] {
	return &LoaderCache[KeyT, ValueT]{baseCache: baseCache}
}

func (c *LoaderCache[KeyT, ValueT]) Set(key KeyT, value ValueT) {
	_ = c.baseCache.Set(key, value)
}

func (c *LoaderCache[KeyT, ValueT]) Get(key KeyT) (ValueT, bool) {
	return c.baseCache.Get(key)
}

func (c *LoaderCache[KeyT, ValueT]) Del(key KeyT) {
	_ = c.baseCache.Del(key)
}
