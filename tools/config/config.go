// tools/config/config.go
package config

// ConfigManager 配置管理器
type ConfigManager struct {
    configs   map[string]*Config
    watchers  map[string][]ConfigWatcher
    loader    ConfigLoader
    validator ConfigValidator
    cache     *cache.Cache[string, interface{}]
    mu        sync.RWMutex
}

// Config 配置结构
type Config struct {
    Data       interface{}
    Format     ConfigFormat
    Source     ConfigSource
    LastUpdate time.Time
    Version    int64
}

// ConfigSource 配置源接口
type ConfigSource interface {
    Load() ([]byte, error)
    Save([]byte) error
    Watch() (<-chan ConfigEvent, error)
}

// ConfigValidator 配置验证器
type ConfigValidator interface {
    Validate(interface{}) error
    Schema() interface{}
}
