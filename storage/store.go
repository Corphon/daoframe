// storage/store.go
package storage

import (
    "context"
    "sync"
    
    "github.com/Corphon/daoframe/tools/cache"
    "github.com/Corphon/daoframe/tools/metrics"
)

// Store 存储接口
type Store interface {
    // 基本操作
    Get(ctx context.Context, key string) (*Item, error)
    Set(ctx context.Context, key string, value []byte, opts *Options) error
    Delete(ctx context.Context, key string) error
    
    // 批量操作
    BatchGet(ctx context.Context, keys []string) (map[string]*Item, error)
    BatchSet(ctx context.Context, items map[string]*Item) error
    BatchDelete(ctx context.Context, keys []string) error
    
    // 查询操作
    List(ctx context.Context, filter *Filter) ([]*Item, error)
    Count(ctx context.Context, filter *Filter) (int64, error)
    
    // 事务操作
    Begin(ctx context.Context) (Transaction, error)
    
    // 维护操作
    Cleanup(ctx context.Context) error
    Compact(ctx context.Context) error
    Backup(ctx context.Context, dst string) error
    
    // 监控操作
    Stats() *StoreStats
    Health() *HealthStatus
}

// StoreStats 存储统计
type StoreStats struct {
    ItemCount     int64
    StorageSize   int64
    OperationStats OperationStats
    CacheStats    *cache.Stats
}

// OperationStats 操作统计
type OperationStats struct {
    Reads         uint64
    Writes        uint64
    Deletes       uint64
    CacheHits     uint64
    CacheMisses   uint64
    Errors        uint64
}

// BaseStore 基础存储实现
type BaseStore struct {
    mu          sync.RWMutex
    items       map[string]*Item
    cache       *cache.Cache
    metrics     *StoreMetrics
    compressor  Compressor
    encryptor   Encryptor
    logger      Logger
}

func NewBaseStore(opts ...StoreOption) *BaseStore {
    store := &BaseStore{
        items:  make(map[string]*Item),
        cache:  cache.New(cache.DefaultExpiration),
        metrics: newStoreMetrics(),
    }
    
    // 应用选项
    for _, opt := range opts {
        opt(store)
    }
    
    return store
}

// Get 获取数据
func (s *BaseStore) Get(ctx context.Context, key string) (*Item, error) {
    // 检查缓存
    if item, found := s.cache.Get(key); found {
        s.metrics.CacheHits.Inc()
        return item.(*Item), nil
    }
    s.metrics.CacheMisses.Inc()
    
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    item, exists := s.items[key]
    if !exists {
        return nil, errors.NotFound("key not found: %s", key)
    }
    
    // 检查过期
    if !item.ExpireAt.IsZero() && item.ExpireAt.Before(time.Now()) {
        return nil, errors.NotFound("key expired: %s", key)
    }
    
    // 解密
    if s.encryptor != nil {
        value, err := s.encryptor.Decrypt(item.Value)
        if err != nil {
            return nil, err
        }
        item.Value = value
    }
    
    // 解压
    if s.compressor != nil {
        value, err := s.compressor.Decompress(item.Value)
        if err != nil {
            return nil, err
        }
        item.Value = value
    }
    
    // 更新缓存
    s.cache.Set(key, item, cache.DefaultExpiration)
    
    s.metrics.Reads.Inc()
    return item, nil
}
