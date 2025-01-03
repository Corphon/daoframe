//system/types.go
package system

import (
    "time"
    
    "github.com/Corphon/daoframe/core"
    "github.com/Corphon/daoframe/model"
)

// SystemType 系统类型
type SystemType uint8

const (
    UniverseSystem SystemType = iota
    InteractionSystem
    EvolutionSystem
    MonitorSystem
)

// SystemState 系统状态
type SystemState uint8

const (
    SystemStateInactive SystemState = iota
    SystemStateStarting
    SystemStateRunning
    SystemStatePaused
    SystemStateStopping
    SystemStateStopped
)

// SystemConfig 系统配置
type SystemConfig struct {
    UpdateInterval  time.Duration
    BufferSize     int
    MaxWorkers     int
    EnableMetrics  bool
}
