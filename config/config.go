// config/config.go

package config

import (
    "time"
    "encoding/json"
    "github.com/Corphon/daoframe/errors"
)

// CoreConfig 框架核心配置
type CoreConfig struct {
    // 基础配置
    Debug          bool          `json:"debug"`
    LogLevel       string        `json:"log_level"`
    MaxGoroutines  int          `json:"max_goroutines"`
    
    // 超时设置
    DefaultTimeout time.Duration `json:"default_timeout"`
    MaxTimeout     time.Duration `json:"max_timeout"`
    
    // 监控设置
    MetricsEnabled bool          `json:"metrics_enabled"`
    MetricsPort    int          `json:"metrics_port"`
    
    // 生命周期配置
    LifeCycleConfig `json:"lifecycle"`
    
    // 调度器配置
    SchedulerConfig `json:"scheduler"`
}

// LifeCycleConfig 生命周期管理配置
type LifeCycleConfig struct {
    // 清理配置
    CleanupInterval   time.Duration `json:"cleanup_interval"`
    MaxInactiveTime   time.Duration `json:"max_inactive_time"`
    
    // 分片设置
    ShardCount        int          `json:"shard_count"`
    
    // 观察者设置
    EnableObservers   bool         `json:"enable_observers"`
    ObserverBuffer    int          `json:"observer_buffer"`
}

// SchedulerConfig 调度器配置
type SchedulerConfig struct {
    // 任务配置
    MaxConcurrentTasks int           `json:"max_concurrent_tasks"`
    TaskQueueSize      int           `json:"task_queue_size"`
    
    // 调度配置
    ScheduleInterval   time.Duration `json:"schedule_interval"`
    RetryAttempts     int           `json:"retry_attempts"`
    RetryDelay        time.Duration `json:"retry_delay"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *CoreConfig {
    return &CoreConfig{
        Debug:          false,
        LogLevel:       "info",
        MaxGoroutines:  10000,
        DefaultTimeout: time.Second * 30,
        MaxTimeout:     time.Minute * 5,
        MetricsEnabled: true,
        MetricsPort:    9090,
        
        LifeCycleConfig: LifeCycleConfig{
            CleanupInterval: time.Hour,
            MaxInactiveTime: time.Hour * 24,
            ShardCount:      32,
            EnableObservers: true,
            ObserverBuffer:  1000,
        },
        
        SchedulerConfig: SchedulerConfig{
            MaxConcurrentTasks: 100,
            TaskQueueSize:      1000,
            ScheduleInterval:   time.Second,
            RetryAttempts:     3,
            RetryDelay:        time.Second * 5,
        },
    }
}

// Validate 验证配置
func (c *CoreConfig) Validate() error {
    // 验证基础配置
    if c.MaxGoroutines <= 0 {
        return errors.New("max_goroutines must be positive")
    }
    
    // 验证超时设置
    if c.DefaultTimeout <= 0 {
        return errors.New("default_timeout must be positive")
    }
    if c.MaxTimeout < c.DefaultTimeout {
        return errors.New("max_timeout must be greater than default_timeout")
    }
    
    // 验证生命周期配置
    if err := c.validateLifeCycleConfig(); err != nil {
        return err
    }
    
    // 验证调度器配置
    if err := c.validateSchedulerConfig(); err != nil {
        return err
    }
    
    return nil
}

// validateLifeCycleConfig 验证生命周期配置
func (c *CoreConfig) validateLifeCycleConfig() error {
    if c.CleanupInterval <= 0 {
        return errors.New("cleanup_interval must be positive")
    }
    if c.MaxInactiveTime <= c.CleanupInterval {
        return errors.New("max_inactive_time must be greater than cleanup_interval")
    }
    if c.ShardCount <= 0 {
        return errors.New("shard_count must be positive")
    }
    if c.EnableObservers && c.ObserverBuffer <= 0 {
        return errors.New("observer_buffer must be positive when observers are enabled")
    }
    return nil
}

// validateSchedulerConfig 验证调度器配置
func (c *CoreConfig) validateSchedulerConfig() error {
    if c.MaxConcurrentTasks <= 0 {
        return errors.New("max_concurrent_tasks must be positive")
    }
    if c.TaskQueueSize <= 0 {
        return errors.New("task_queue_size must be positive")
    }
    if c.ScheduleInterval <= 0 {
        return errors.New("schedule_interval must be positive")
    }
    if c.RetryAttempts < 0 {
        return errors.New("retry_attempts cannot be negative")
    }
    if c.RetryAttempts > 0 && c.RetryDelay <= 0 {
        return errors.New("retry_delay must be positive when retries are enabled")
    }
    return nil
}

// LoadFromFile 从文件加载配置
func LoadFromFile(filename string) (*CoreConfig, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    
    config := DefaultConfig()
    if err := json.Unmarshal(data, config); err != nil {
        return nil, err
    }
    
    if err := config.Validate(); err != nil {
        return nil, err
    }
    
    return config, nil
}

// LoadFromEnv 从环境变量加载配置
func LoadFromEnv() (*CoreConfig, error) {
    config := DefaultConfig()
    
    // 从环境变量加载配置
    if debug := os.Getenv("DAOFRAME_DEBUG"); debug != "" {
        config.Debug = debug == "true"
    }
    
    if logLevel := os.Getenv("DAOFRAME_LOG_LEVEL"); logLevel != "" {
        config.LogLevel = logLevel
    }
    
    // ... 其他环境变量加载逻辑
    
    if err := config.Validate(); err != nil {
        return nil, err
    }
    
    return config, nil
}
