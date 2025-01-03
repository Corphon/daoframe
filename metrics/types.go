// metrics/types.go
package metrics

import (
    "sync/atomic"
    "time"
)

// MetricType 指标类型
type MetricType string

const (
    CounterMetric   MetricType = "counter"   // 计数器类型
    GaugeMetric     MetricType = "gauge"     // 仪表类型
    HistogramMetric MetricType = "histogram" // 直方图类型
    SummaryMetric   MetricType = "summary"   // 摘要类型
)

// MetricValue 指标值接口
type MetricValue interface {
    Add(float64)
    Set(float64)
    Get() float64
    Reset()
}

// Counter 计数器实现
type Counter struct {
    value uint64
}

func (c *Counter) Add(delta float64) {
    if delta < 0 {
        return // 计数器只能增加
    }
    atomic.AddUint64(&c.value, uint64(delta))
}

func (c *Counter) Set(v float64) {
    atomic.StoreUint64(&c.value, uint64(v))
}

func (c *Counter) Get() float64 {
    return float64(atomic.LoadUint64(&c.value))
}

func (c *Counter) Reset() {
    atomic.StoreUint64(&c.value, 0)
}

// Gauge 仪表实现
type Gauge struct {
    value uint64
}

func (g *Gauge) Add(delta float64) {
    var new uint64
    for {
        old := atomic.LoadUint64(&g.value)
        new = math.Float64bits(math.Float64frombits(old) + delta)
        if atomic.CompareAndSwapUint64(&g.value, old, new) {
            break
        }
    }
}

// Histogram 直方图实现
type Histogram struct {
    mutex      sync.RWMutex
    count      uint64
    sum        float64
    buckets    []float64
    counts     []uint64
    quantiles  []float64
}

// HistogramOpts 直方图配置
type HistogramOpts struct {
    Buckets   []float64
    Quantiles []float64
    BufSize   int
}

func NewHistogram(opts HistogramOpts) *Histogram {
    if len(opts.Buckets) == 0 {
        opts.Buckets = DefaultBuckets
    }
    sort.Float64s(opts.Buckets)
    
    return &Histogram{
        buckets:   opts.Buckets,
        counts:    make([]uint64, len(opts.Buckets)),
        quantiles: opts.Quantiles,
    }
}
