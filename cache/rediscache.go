package cache

import (
	"strconv"
)

// RedisCache 待实现
type RedisCache[V any] struct {
	Name      string
	PrefixKey string
}

func NewRedisCache[V any](name string) *RedisCache[V] {
	//根据泛型V生成prefixKey
	prefixKey := getPrefixKey(*new(V))
	return &RedisCache[V]{Name: name, PrefixKey: prefixKey}
}

func getPrefixKey(a any) string {
	prefixKey := hash32(a)
	return strconv.Itoa(int(prefixKey))
}

func (cache *RedisCache[V]) GetName() string {
	return cache.Name
}

func (cache *RedisCache[V]) GetPrefixKey() string {
	return cache.PrefixKey
}

func (cache *RedisCache[V]) Get(k string) []*V {
	panic("not implemented")
	return nil
}

func (cache *RedisCache[V]) Put(k string, v []*V) {
	panic("not implemented")
}

func (cache *RedisCache[V]) Remove(k string) {
	panic("not implemented")
}

func (cache *RedisCache[V]) Clear() {
	//todo 这里的clear清除的就是redis中key的前缀为cache.PrefixKey的键值对
	panic("not implemented")
}
