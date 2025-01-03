//config/validator.go
package config

// CoreConfigValidator 核心配置验证器
type CoreConfigValidator struct{}

func (v *CoreConfigValidator) Validate(cfg interface{}) error {
    config, ok := cfg.(*CoreConfig)
    if !ok {
        return errors.New("invalid config type")
    }
    
    // 验证基础配置
    if err := v.validateBasicConfig(config); err != nil {
        return err
    }
    
    // 验证生命周期配置
    if err := v.validateLifecycleConfig(&config.LifeCycleConfig); err != nil {
        return err
    }
    
    // 验证调度器配置
    if err := v.validateSchedulerConfig(&config.SchedulerConfig); err != nil {
        return err
    }
    
    return nil
}

func (v *CoreConfigValidator) validateBasicConfig(config *CoreConfig) error {
    if config.MaxGoroutines <= 0 {
        return errors.New("max_goroutines must be positive")
    }
    if config.DefaultTimeout <= 0 {
        return errors.New("default_timeout must be positive")
    }
    if config.MaxTimeout < config.DefaultTimeout {
        return errors.New("max_timeout must be greater than default_timeout")
    }
    return nil
}
