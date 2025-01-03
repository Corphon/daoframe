// discovery/registry.go
package discovery

import (
    "context"
    "sync"
    "time"
    
    "github.com/Corphon/daoframe/tools/cache"
    "github.com/Corphon/daoframe/errors"
)

// ServiceRegistry 服务注册中心
type ServiceRegistry struct {
    mu          sync.RWMutex
    instances   map[string]*ServiceInstance
    subscribers map[string][]chan *ServiceEvent
    store       RegistryStore
    cache       *cache.Cache
    health      *HealthChecker
    metrics     *RegistryMetrics
    
    ctx         context.Context
    cancel      context.CancelFunc
    wg          sync.WaitGroup
}

// RegistryConfig 注册中心配置
type RegistryConfig struct {
    // 心跳超时时间
    HeartbeatTimeout time.Duration
    // 清理间隔
    CleanupInterval time.Duration
    // 缓存过期时间
    CacheTTL        time.Duration
    // 健康检查配置
    HealthCheck     *HealthCheckConfig
}

// Register 注册服务实例
func (sr *ServiceRegistry) Register(ctx context.Context, instance *ServiceInstance) error {
    if err := sr.validateInstance(instance); err != nil {
        return err
    }

    sr.mu.Lock()
    defer sr.mu.Unlock()

    // 检查是否已存在
    if existing, exists := sr.instances[instance.ID]; exists {
        return errors.AlreadyExists("service instance already registered: %s", instance.ID)
    }

    // 设置注册时间
    instance.RegisterTime = time.Now()
    instance.LastHeartbeat = time.Now()

    // 保存到存储
    if err := sr.store.SaveInstance(ctx, instance); err != nil {
        return err
    }

    // 更新内存状态
    sr.instances[instance.ID] = instance

    // 发送事件
    sr.publishEvent(&ServiceEvent{
        Type:      EventRegister,
        Instance:  instance,
        Timestamp: time.Now(),
    })

    // 启动健康检查
    sr.health.StartCheck(instance)

    sr.metrics.RegisteredInstances.Inc()
    return nil
}

// Deregister 注销服务实例
func (sr *ServiceRegistry) Deregister(ctx context.Context, instanceID string) error {
    sr.mu.Lock()
    defer sr.mu.Unlock()

    instance, exists := sr.instances[instanceID]
    if !exists {
        return errors.NotFound("service instance not found: %s", instanceID)
    }

    // 从存储中删除
    if err := sr.store.DeleteInstance(ctx, instanceID); err != nil {
        return err
    }

    // 停止健康检查
    sr.health.StopCheck(instanceID)

    // 删除内存状态
    delete(sr.instances, instanceID)

    // 发送事件
    sr.publishEvent(&ServiceEvent{
        Type:      EventDeregister,
        Instance:  instance,
        Timestamp: time.Now(),
    })

    sr.metrics.RegisteredInstances.Dec()
    return nil
}
