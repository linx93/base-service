package cache

import (
	"github.com/ohmountain/memcache"
)

// DefaultCache 默认缓存，基于山哥写的缓存，其实就是包装了一下山哥的患缓存库
type DefaultCache[V any] struct {
	Name  string
	Store *memcache.Memcache[[]*V]
}

func (cache *DefaultCache[V]) GetName() string {
	return cache.Name
}

func (cache *DefaultCache[V]) Get(k string) []*V {
	if ext := cache.Store.Get(k); ext != nil {
		return *ext
	}
	return nil
}

func (cache *DefaultCache[V]) Put(k string, v []*V) {
	cache.Store.Set(k, v)
}

func (cache *DefaultCache[V]) Remove(k string) {
	cache.Store.Delete(k)
}

func (cache *DefaultCache[V]) Clear() {
	cache.Store.Clear()
}
