//storage/manager.go
package storage

import (
    "context"
    "sync"
    
    "github.com/Corphon/daoframe/tools/async"
)

// StorageManager 存储管理器
type StorageManager struct {
    stores     map[string]Store
    factory    *StoreFactory
    monitor    *StoreMonitor
    replicator *StoreReplicator
    
    // 配置
    config     *ManagerConfig
    // 工作池
    pool       *async.WorkerPool
    // 监控
    metrics    *ManagerMetrics
    
    ctx        context.Context
    cancel     context.CancelFunc
    wg         sync.WaitGroup
}

// ManagerConfig 管理器配置
type ManagerConfig struct {
    // 存储配置
    StoreConfigs map[string]*StoreConfig
    // 复制配置
    ReplicationConfig *ReplicationConfig
    // 监控配置
    MonitoringConfig *MonitoringConfig
    // 工作池配置
    PoolConfig *async.PoolConfig
}

// StoreMonitor 存储监控器
type StoreMonitor struct {
    manager    *StorageManager
    collectors []metrics.Collector
    alerts     []Alert
    checks     []HealthCheck
}

// StoreReplicator 存储复制器
type StoreReplicator struct {
    manager    *StorageManager
    strategy   ReplicationStrategy
    syncer     *StoreSyncer
}

// CreateStore 创建存储
func (m *StorageManager) CreateStore(ctx context.Context, name string, typ StoreType) (Store, error) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    // 检查是否已存在
    if _, exists := m.stores[name]; exists {
        return nil, errors.AlreadyExists("store already exists: %s", name)
    }
    
    // 创建存储
    store, err := m.factory.Create(typ)
    if err != nil {
        return nil, err
    }
    
    // 初始化监控
    if err := m.monitor.RegisterStore(store); err != nil {
        return nil, err
    }
    
    // 配置复制
    if m.config.ReplicationConfig.Enabled {
        if err := m.replicator.ConfigureStore(store); err != nil {
            return nil, err
        }
    }
    
    m.stores[name] = store
    return store, nil
}
