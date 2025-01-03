//config/manager.go
package config

import (
    "sync"
    "time"
)

// ConfigManager 配置管理器
type ConfigManager struct {
    mu          sync.RWMutex
    configs     map[ConfigType]interface{}
    validators  map[ConfigType]ConfigValidator
    watchers    map[ConfigType][]ConfigWatcher
    lastUpdate  time.Time
}

// ConfigValidator 配置验证接口
type ConfigValidator interface {
    Validate(interface{}) error
}

// ConfigWatcher 配置变更观察者接口
type ConfigWatcher interface {
    OnConfigChange(ConfigType, interface{})
}

func NewConfigManager() *ConfigManager {
    return &ConfigManager{
        configs:    make(map[ConfigType]interface{}),
        validators: make(map[ConfigType]ConfigValidator),
        watchers:   make(map[ConfigType][]ConfigWatcher),
    }
}

// RegisterConfig 注册配置
func (cm *ConfigManager) RegisterConfig(typ ConfigType, config interface{}, validator ConfigValidator) error {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    if validator != nil {
        if err := validator.Validate(config); err != nil {
            return err
        }
    }
    
    cm.configs[typ] = config
    cm.validators[typ] = validator
    cm.lastUpdate = time.Now()
    
    // 通知观察者
    cm.notifyWatchers(typ, config)
    return nil
}
