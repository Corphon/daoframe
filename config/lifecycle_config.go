// config/lifecycle_config.go

package config

type LifeCycleConfig struct {
    // 清理配置
    CleanupInterval   time.Duration
    MaxInactiveTime   time.Duration
    
    // 分片锁配置
    ShardCount        int
    
    // 观察者配置
    EnableObservers   bool
}

func DefaultLifeCycleConfig() *LifeCycleConfig {
    return &LifeCycleConfig{
        CleanupInterval: time.Hour,
        MaxInactiveTime: time.Hour * 24,
        ShardCount:      32,
        EnableObservers: true,
    }
}
