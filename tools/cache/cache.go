// tools/cache/cache.go
package cache

// Cache 通用缓存接口
type Cache[K comparable, V any] interface {
    Get(key K) (V, bool)
    Set(key K, value V, ttl time.Duration) error
    Delete(key K)
    Clear()
    Size() int
}

// LRUCache LRU缓存实现
type LRUCache[K comparable, V any] struct {
    capacity int
    items    map[K]*lruItem[K, V]
    head     *lruItem[K, V]
    tail     *lruItem[K, V]
    mu       sync.RWMutex
}

type lruItem[K comparable, V any] struct {
    key       K
    value     V
    expiry    time.Time
    prev, next *lruItem[K, V]
}
