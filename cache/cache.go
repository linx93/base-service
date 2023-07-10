package cache

// ICache 基础缓存抽象
type ICache[V any] interface {
	// GetName 获取缓存名称
	GetName() string

	// Get 获取缓存
	Get(key string) []*V

	// Put 设置缓存
	Put(K string, v []*V)

	// Remove 移除缓存
	Remove(k string)

	// Clear 清除缓存
	Clear()
}
