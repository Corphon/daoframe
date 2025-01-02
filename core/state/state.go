// core/state/state.go

package state

import (
    "fmt"
)

// BaseState 基础状态
type BaseState uint8

const (
    StateVoid BaseState = iota
    StateInactive
    StateActive
    StatePaused
    StateTerminated
)

// LifecycleState 生命周期状态
type LifecycleState uint8

const (
    StageOrigin LifecycleState = iota
    StageBirth
    StageGrowth
    StagePeak
    StageDecline
    StageEnd
    StageReturn
)

// 新增状态管理器
type StateManager struct {
    baseState      BaseState
    lifecycleState LifecycleState
    mu            sync.RWMutex
}

// StateTransitionMap 定义所有可能的状态转换
var StateTransitionMap = map[State][]State{
    // 基础状态转换
    StateVoid:      {StateInactive, StageOrigin},  // 可以转向基础激活态或生命周期起源
    StateInactive:  {StateActive},
    StateActive:    {StatePaused, StateTerminated},
    StatePaused:    {StateActive, StateTerminated},
    StateTerminated: {StateVoid},  // 循环返回虚无

    // 生命周期状态转换（从原 lifecycle_state.go 迁移）
    StageOrigin:    {StageBirth},
    StageBirth:     {StageGrowth},
    StageGrowth:    {StagePeak, StageDecline},
    StagePeak:      {StageDecline},
    StageDecline:   {StageEnd},
    StageEnd:       {StageReturn},
    StageReturn:    {StateVoid},  // 生命周期结束回到虚无
}

// ValidateTransition 验证状态转换是否有效
func ValidateTransition(current, next State) error {
    validNextStates, exists := StateTransitionMap[current]
    if !exists {
        return fmt.Errorf("invalid current state: %v", current)
    }
    
    for _, validState := range validNextStates {
        if next == validState {
            return nil
        }
    }
    return fmt.Errorf("invalid state transition from %v to %v", current, next)
}

// GetStateName 获取状态的字符串表示
func GetStateName(s State) string {
    switch s {
    case StateVoid:
        return "Void"
    case StateInactive:
        return "Inactive"
    case StateActive:
        return "Active"
    case StatePaused:
        return "Paused"
    case StateTerminated:
        return "Terminated"
    case StageOrigin:
        return "Origin"
    case StageBirth:
        return "Birth"
    case StageGrowth:
        return "Growth"
    case StagePeak:
        return "Peak"
    case StageDecline:
        return "Decline"
    case StageEnd:
        return "End"
    case StageReturn:
        return "Return"
    default:
        return "Unknown"
    }
}
