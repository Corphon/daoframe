//plugin/container.go
package plugin

import (
    "context"
    "sync"
)

// PluginContainer 插件容器
type PluginContainer struct {
    plugin     Plugin
    state      *PluginState
    config     *ContainerConfig
    isolator   *PluginIsolator
    monitor    *PluginMonitor
    lifecycle  *PluginLifecycle
    
    // 资源管理
    resources  *ResourceManager
    // 错误处理
    errorHandler ErrorHandler
    // 监控指标
    metrics    *ContainerMetrics
    
    mu         sync.RWMutex
}

// ContainerConfig 容器配置
type ContainerConfig struct {
    // 隔离选项
    IsolationOpts IsolationOptions
    // 资源限制
    ResourceLimits ResourceLimits
    // 监控选项
    MonitorOpts MonitorOptions
    // 生命周期钩子
    Lifecycle LifecycleHooks
}

// PluginIsolator 插件隔离器
type PluginIsolator struct {
    namespace string
    cgroups   *CGroupManager
    network   *NetworkNamespace
    filesystem *FilesystemIsolator
}

// ResourceManager 资源管理器
type ResourceManager struct {
    limits    ResourceLimits
    usage     *ResourceUsage
    allocator *ResourceAllocator
    monitor   *ResourceMonitor
}

// Init 初始化容器
func (pc *PluginContainer) Init(ctx context.Context) error {
    pc.mu.Lock()
    defer pc.mu.Unlock()
    
    // 设置隔离环境
    if err := pc.isolator.Setup(); err != nil {
        return err
    }
    
    // 分配资源
    if err := pc.resources.Allocate(); err != nil {
        return err
    }
    
    // 启动监控
    if err := pc.monitor.Start(ctx); err != nil {
        return err
    }
    
    // 初始化插件
    if err := pc.plugin.Init(ctx); err != nil {
        return err
    }
    
    // 更新状态
    pc.setState(PluginStateActive)
    
    return nil
}
