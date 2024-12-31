package state

// State 统一的状态定义
type State uint8

const (
    StateVoid      State = iota // 虚无状态
    StateInactive              // 未激活
    StateActive               // 激活
    StatePaused              // 暂停
    StateTerminated          // 终止
    
    // 生命周期特有状态
    StateOrigin               // 起源
    StageBirth               // 诞生
    StageGrowth              // 成长
    StagePeak                // 巅峰
    StageDecline             // 衰退
    StageEnd                 // 终末
    StageReturn              // 返归
)

// StateTransitionMap 定义状态转换规则
var StateTransitionMap = map[State][]State{
    StateVoid:    {StateInactive, StateOrigin},
    StateInactive: {StateActive},
    StateActive:   {StatePaused, StateTerminated},
    StatePaused:   {StateActive, StateTerminated},
    StateTerminated: {StateVoid},
    // 生命周期状态转换
    StateOrigin:  {StageBirth},
    StageBirth:   {StageGrowth},
    StageGrowth:  {StagePeak, StageDecline},
    StagePeak:    {StageDecline},
    StageDecline: {StageEnd},
    StageEnd:     {StageReturn},
    StageReturn:  {StateVoid},
}

// ValidateTransition 统一的状态转换验证
func ValidateTransition(current, next State) bool {
    validStates, exists := StateTransitionMap[current]
    if !exists {
        return false
    }
    
    for _, validState := range validStates {
        if next == validState {
            return true
        }
    }
    return false
}
