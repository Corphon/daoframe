// plugin/plugin.go
package plugin

import (
    "context"
    "sync/atomic"
)

// Plugin 插件接口
type Plugin interface {
    // 基本信息
    Info() *PluginInfo
    State() PluginState
    
    // 生命周期管理
    Init(ctx context.Context) error
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Reload(ctx context.Context) error
    
    // 健康检查
    HealthCheck(ctx context.Context) error
    
    // 资源管理
    Cleanup() error
    
    // 扩展点
    GetExtensionPoints() []ExtensionPoint
    GetExtensions() []Extension
}

// PluginImpl 插件基础实现
type PluginImpl struct {
    info      *PluginInfo
    state     atomic.Value
    container *PluginContainer
    config    *PluginConfig
    logger    Logger
    metrics   *PluginMetrics
    
    extensions       map[string]Extension
    extensionPoints map[string]ExtensionPoint
    
    ctx     context.Context
    cancel  context.CancelFunc
    running atomic.Bool
}

// NewPlugin 创建新插件
func NewPlugin(info *PluginInfo, opts ...PluginOption) (*PluginImpl, error) {
    if err := validatePluginInfo(info); err != nil {
        return nil, err
    }
    
    ctx, cancel := context.WithCancel(context.Background())
    
    p := &PluginImpl{
        info:            info,
        extensions:      make(map[string]Extension),
        extensionPoints: make(map[string]ExtensionPoint),
        ctx:            ctx,
        cancel:         cancel,
    }
    
    // 应用选项
    for _, opt := range opts {
        if err := opt(p); err != nil {
            cancel()
            return nil, err
        }
    }
    
    // 初始化指标
    p.metrics = NewPluginMetrics(info.ID)
    
    // 设置初始状态
    p.state.Store(PluginStateInitializing)
    
    return p, nil
}

// RegisterExtension 注册扩展
func (p *PluginImpl) RegisterExtension(ext Extension) error {
    if ext == nil {
        return errors.New("extension cannot be nil")
    }
    
    p.mu.Lock()
    defer p.mu.Unlock()
    
    id := ext.ID()
    if _, exists := p.extensions[id]; exists {
        return errors.New("extension already registered")
    }
    
    p.extensions[id] = ext
    p.metrics.ExtensionsRegistered.Inc()
    
    return nil
}
