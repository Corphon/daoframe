package basic

type Trigram uint8

const (
    TrigramQian Trigram = iota // 乾 ☰
    TrigramKun                 // 坤 ☷
    TrigramZhen                // 震 ☳
    TrigramXun                 // 巽 ☴
    TrigramKan                 // 坎 ☵
    TrigramLi                  // 离 ☲
    TrigramGen                 // 艮 ☶
    TrigramDui                 // 兑 ☱
)

// BaGua 八卦系统
type BaGua struct {
    mu            sync.RWMutex
    trigrams      map[Trigram]*TrigramState
    energyFlows   map[Trigram][]EnergyFlow
    interactions  map[Trigram]map[Trigram]float64
    ctx           *core.DaoContext
}

// TrigramState 卦象状态
type TrigramState struct {
    trigram    Trigram
    direction  Direction
    energy     float64
    attribute  *YinYangAttribute
    element    Phase
}

// EnergyFlow 能量流动
type EnergyFlow struct {
    source    Trigram
    target    Trigram
    strength  float64
    nature    Nature
}
