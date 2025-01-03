//discovery/resolver.go
package discovery

import (
    "context"
    "sync"
    "time"
)

// ServiceResolver 服务解析器
type ServiceResolver struct {
    registry    *ServiceRegistry
    balancer    *LoadBalancer
    cache       *ServiceCache
    watcher     *ServiceWatcher
    metrics     *ResolverMetrics
    
    mu          sync.RWMutex
    watchers    map[string]*instanceWatcher
}

// ResolveOptions 解析选项
type ResolveOptions struct {
    Strategy   string            // 负载均衡策略
    Filters    []InstanceFilter  // 实例过滤器
    Timeout    time.Duration     // 解析超时
    Cache      bool              // 是否使用缓存
    Retry      int              // 重试次数
}

// Resolve 解析服务实例
func (sr *ServiceResolver) Resolve(ctx context.Context, serviceName string, opts *ResolveOptions) (*ServiceInstance, error) {
    // 检查缓存
    if opts.Cache {
        if instance := sr.cache.Get(serviceName); instance != nil {
            sr.metrics.CacheHits.Inc()
            return instance, nil
        }
        sr.metrics.CacheMisses.Inc()
    }

    // 获取服务实例列表
    instances, err := sr.registry.ListInstances(ctx, serviceName)
    if err != nil {
        return nil, err
    }

    // 应用过滤器
    filtered := sr.applyFilters(instances, opts.Filters)
    if len(filtered) == 0 {
        return nil, errors.NoAvailableInstances
    }

    // 负载均衡选择实例
    instance, err := sr.balancer.Select(filtered, opts.Strategy)
    if err != nil {
        return nil, err
    }

    // 更新缓存
    if opts.Cache {
        sr.cache.Set(serviceName, instance)
    }

    sr.metrics.ResolutionsTotal.Inc()
    return instance, nil
}
