package model

type ElementEvent struct {
    EntityID    string
    ElementType Phase
    OldStrength uint8
    NewStrength uint8
    Timestamp   time.Time
}

type BalanceEvent struct {
    EntityID     string
    BalanceType  string
    Measurements map[string]float64
    Timestamp    time.Time
}

type CycleEvent struct {
    EntityID   string
    CycleType  string
    Duration   time.Duration
    Changes    []StateChange
    Timestamp  time.Time
}
