//config/watcher.go
package config

// ConfigWatcherFunc 配置观察者函数类型
type ConfigWatcherFunc func(ConfigType, interface{})

// ConfigChangeEvent 配置变更事件
type ConfigChangeEvent struct {
    Type     ConfigType
    OldValue interface{}
    NewValue interface{}
    Time     time.Time
}

// ConfigWatcherHub 配置观察者中心
type ConfigWatcherHub struct {
    mu       sync.RWMutex
    watchers map[ConfigType][]ConfigWatcher
    events   chan ConfigChangeEvent
}

func NewConfigWatcherHub() *ConfigWatcherHub {
    return &ConfigWatcherHub{
        watchers: make(map[ConfigType][]ConfigWatcher),
        events:   make(chan ConfigChangeEvent, 100),
    }
}

// AddWatcher 添加观察者
func (h *ConfigWatcherHub) AddWatcher(typ ConfigType, watcher ConfigWatcher) {
    h.mu.Lock()
    defer h.mu.Unlock()
    
    if h.watchers[typ] == nil {
        h.watchers[typ] = make([]ConfigWatcher, 0)
    }
    h.watchers[typ] = append(h.watchers[typ], watcher)
}
