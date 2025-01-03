//model/types.go
package model

// 统一定义所有基础类型
type (
    Nature uint8
    Phase  uint8
    Force  uint8
    ID     string
)

const (
    // Nature 类型常量
    NatureYin  Nature = iota
    NatureYang
    NatureTai

    // Phase 类型常量
    PhaseWood Phase = iota
    PhaseFire
    PhaseEarth
    PhaseMetal
    PhaseWater
)

// Element 统一的元素接口
type Element interface {
    GetPhase() Phase
    GetNature() Nature
    GetStrength() uint8
    SetStrength(uint8) error
}
