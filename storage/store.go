// storage/store.go
package storage

type Store interface {
    Get(key string) ([]byte, error)
    Set(key string, value []byte) error
    Delete(key string) error
    List(prefix string) ([]string, error)
}

// storage/manager.go
type StorageManager struct {
    stores    map[string]Store
    cache     *cache.Cache
    indexer   *Indexer
    metrics   *StorageMetrics
}
