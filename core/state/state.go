// core/state/state.go

package state

import (
    "fmt"
)

// State 统一的状态定义
type State uint8

const (
    // 基础状态
    StateVoid      State = iota  // 虚无状态
    StateInactive               // 未激活
    StateActive                // 激活
    StatePaused               // 暂停
    StateTerminated          // 终止

    // 生命周期状态
    StageOrigin               // 起源
    StageBirth               // 诞生
    StageGrowth              // 成长
    StagePeak                // 巅峰
    StageDecline             // 衰退
    StageEnd                 // 终末
    StageReturn              // 返归
)

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
