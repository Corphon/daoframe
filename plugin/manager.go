//plugin/manager.go
package plugin

import (
    "context"
    "sync"
    "github.com/Corphon/daoframe/tools/async"
)

// PluginManager 插件管理器
type PluginManager struct {
    mu          sync.RWMutex
    plugins     map[string]*PluginContainer
    loader      *PluginLoader
    registry    *PluginRegistry
    hooks       map[string][]PluginHook
    depResolver *DependencyResolver
    
    // 配置
    config      *ManagerConfig
    // 工作池
    pool        *async.WorkerPool
    // 事件总线
    eventBus    EventBus
    // 监控
    metrics     *ManagerMetrics
    
    ctx         context.Context
    cancel      context.CancelFunc
    wg          sync.WaitGroup
}

// ManagerConfig 管理器配置
type ManagerConfig struct {
    // 插件目录
    PluginDirs []string
    // 并发限制
    MaxConcurrent int
    // 加载超时
    LoadTimeout time.Duration
    // 自动重载
    AutoReload bool
    // 依赖检查
    CheckDependencies bool
    // 版本约束
    VersionConstraints map[string]string
}

// LoadPlugin 加载插件
func (pm *PluginManager) LoadPlugin(ctx context.Context, path string) error {
    // 创建加载任务
    task := &PluginLoadTask{
        Path:    path,
        Context: ctx,
    }
    
    // 提交到工作池
    err := pm.pool.Submit(func() error {
        // 检查插件是否已加载
        if pm.isPluginLoaded(path) {
            return errors.New("plugin already loaded")
        }
        
        // 加载插件
        plugin, err := pm.loader.Load(task)
        if err != nil {
            pm.metrics.LoadFailures.Inc()
            return err
        }
        
        // 检查依赖
        if pm.config.CheckDependencies {
            if err := pm.checkDependencies(plugin); err != nil {
                return err
            }
        }
        
        // 创建容器
        container := NewPluginContainer(plugin, pm.config)
        
        // 注册插件
        if err := pm.registerPlugin(container); err != nil {
            return err
        }
        
        // 初始化插件
        if err := container.Init(ctx); err != nil {
            pm.unregisterPlugin(plugin.Info().ID)
            return err
        }
        
        // 启动插件
        if err := container.Start(ctx); err != nil {
            pm.unregisterPlugin(plugin.Info().ID)
            return err
        }
        
        pm.metrics.LoadedPlugins.Inc()
        return nil
    })
    
    return err
}

// UnloadPlugin 卸载插件
func (pm *PluginManager) UnloadPlugin(ctx context.Context, id string) error {
    pm.mu.Lock()
    defer pm.mu.Unlock()
    
    container, exists := pm.plugins[id]
    if !exists {
        return errors.New("plugin not found")
    }
    
    // 检查依赖关系
    if err := pm.checkDependents(id); err != nil {
        return err
    }
    
    // 停止插件
    if err := container.Stop(ctx); err != nil {
        return err
    }
    
    // 清理资源
    if err := container.Cleanup(); err != nil {
        return err
    }
    
    // 从注册表移除
    delete(pm.plugins, id)
    pm.metrics.LoadedPlugins.Dec()
    
    return nil
}
