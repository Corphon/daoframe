package model

type ModelInteraction struct {
    SourceModel string
    TargetModel string
    Effect      InteractionEffect
    Strength    float64
}

type InteractionManager struct {
    interactions map[string]*ModelInteraction
    observers    []InteractionObserver
    mu          sync.RWMutex
}
