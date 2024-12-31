package basic

// TimeSystem 时序系统
type TimeSystem struct {
    mu         sync.RWMutex
    cycle      *CosmicCycle
    tianGan    *TianGan
    diZhi      *DiZhi
    bagua      *BaGua
    wuXing     *WuXing
}

// CosmicCycle 宇宙周期
type CosmicCycle struct {
    current    CyclePhase
    duration   time.Duration
    patterns   map[CyclePhase]*CyclePattern
}

// CyclePattern 周期模式
type CyclePattern struct {
    phase      CyclePhase
    ganZhi     GanZhiPair
    trigram    Trigram
    element    Phase
    strength   float64
}

// GanZhiPair 天干地支配对
type GanZhiPair struct {
    gan        Gan
    zhi        Zhi
    nature     Nature
    element    Phase
}
