//system/monitor.go
package system

type Monitor struct {
    mu          sync.RWMutex
    collectors  map[string]*MetricCollector
    aggregator  *MetricAggregator
    storage     MetricStorage
    alerts      *AlertManager
    
    config      *MonitorConfig
    state       SystemState
    done        chan struct{}
}

type MetricCollector struct {
    Type       string
    Interval   time.Duration
    Callback   func() []Metric
    Buffer     []Metric
}

type MetricStorage interface {
    Store(metrics []Metric) error
    Query(query MetricQuery) ([]Metric, error)
    Purge(age time.Duration) error
}

type AlertManager struct {
    rules     []AlertRule
    handlers  map[string]AlertHandler
    history   []Alert
}

func NewMonitor(config *MonitorConfig) *Monitor {
    return &Monitor{
        collectors: make(map[string]*MetricCollector),
        aggregator: NewMetricAggregator(),
        alerts:    NewAlertManager(),
        config:    config,
        done:      make(chan struct{}),
    }
}
