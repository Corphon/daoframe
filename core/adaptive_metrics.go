package core

type AdaptiveMetrics struct {
    Effectiveness float64
    Stability     float64
    Balance       float64
    Energy        float64
}

type AdaptiveAction struct {
    Timestamp    time.Time
    ActionType   string
    Effectiveness float64
    Impact       map[string]float64
}
