//storage/types.go
package storage

import (
    "context"
    "time"
    "github.com/Corphon/daoframe/errors"
)

// StoreType 存储类型
type StoreType string

const (
    MemoryStore  StoreType = "memory"
    FileStore    StoreType = "file"
    RedisStore   StoreType = "redis"
    SQLStore     StoreType = "sql"
    MongoStore   StoreType = "mongo"
    EtcdStore    StoreType = "etcd"
)

// Item 存储项
type Item struct {
    Key       string
    Value     []byte
    Metadata  map[string]string
    Version   int64
    Created   time.Time
    Modified  time.Time
    ExpireAt  time.Time
}

// Options 存储选项
type Options struct {
    TTL       time.Duration
    Versioned bool
    Compress  bool
    Encrypt   bool
    Replicate bool
}

// Filter 查询过滤器
type Filter struct {
    Prefix    string
    Tags      map[string]string
    From      time.Time
    To        time.Time
    Limit     int
    Offset    int
    OrderBy   string
    OrderDesc bool
}
