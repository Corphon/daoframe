// model/lifecycle_state.go

package model

// 状态转换定义
var validStateTransitions = map[LifeStage][]LifeStage{
    StageVoid:    {StageOrigin},
    StageOrigin:  {StageBirth},
    StageBirth:   {StageGrowth},
    StageGrowth:  {StagePeak, StageDecline},
    StagePeak:    {StageDecline},
    StageDecline: {StageEnd},
    StageEnd:     {StageReturn},
    StageReturn:  {StageVoid},
}

// 验证状态转换的辅助函数
func ValidateStateTransition(current, next LifeStage) error {
    validNextStages, exists := validStateTransitions[current]
    if !exists {
        return NewLifeCycleError(ErrCodeInvalidState, "Invalid current state")
    }
    
    for _, validStage := range validNextStages {
        if next == validStage {
            return nil
        }
    }
    return NewLifeCycleError(ErrCodeStateTransition, "Invalid state transition")
}
