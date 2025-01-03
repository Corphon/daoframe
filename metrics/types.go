// metrics/types.go
package metrics

type MetricType string

const (
    CounterMetric  MetricType = "counter"
    GaugeMetric    MetricType = "gauge"
    HistogramMetric MetricType = "histogram"
    SummaryMetric   MetricType = "summary"
)

// Metric 指标接口
type Metric interface {
    Name() string
    Type() MetricType
    Labels() map[string]string
    Value() float64
}

// metrics/collector.go
type MetricCollector struct {
    metrics    map[string]Metric
    aggregator *MetricAggregator
    reporter   *MetricReporter
    storage    MetricStorage
}

// metrics/reporter.go
type MetricReporter struct {
    exporters []MetricExporter
    interval  time.Duration
    buffer    *ring.Buffer
}
