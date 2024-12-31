package model

type ElementEffect struct {
    SourcePhase Phase
    TargetPhase Phase
    Strength    float64
    Duration    time.Duration
    Type        EffectType
}

type EffectCalculator struct {
    mutualEffects    map[Phase]map[Phase]float64
    cycleStrength    map[Relationship]float64
    environmentFactor float64
}
