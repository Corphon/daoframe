//metrics/aggregator.go
package metrics

import (
    "sync"
    "time"
)

// MetricAggregator 指标聚合器
type MetricAggregator struct {
    mu         sync.RWMutex
    windows    map[string]*TimeWindow
    resolution time.Duration
    
    // 聚合规则
    rules      map[string]AggregateRule
    // 数据窗口
    ringBuffer *TimeRingBuffer
    // 聚合结果缓存
    cache      *AggregateCache
}

// TimeWindow 时间窗口
type TimeWindow struct {
    start    time.Time
    duration time.Duration
    buckets  []*MetricBucket
}

// MetricBucket 指标桶
type MetricBucket struct {
    count    uint64
    sum      float64
    min      float64
    max      float64
    values   []float64
}

// AggregateRule 聚合规则
type AggregateRule struct {
    Name       string
    Type       MetricType
    Operation  AggregateOperation
    Window     time.Duration
    Quantiles  []float64
}

func (ma *MetricAggregator) Aggregate(metric Metric) (*AggregateResult, error) {
    ma.mu.RLock()
    window, exists := ma.windows[metric.Name()]
    ma.mu.RUnlock()
    
    if !exists {
        return nil, errors.New("no window found for metric")
    }
    
    // 获取聚合规则
    rule, err := ma.getRule(metric.Name())
    if err != nil {
        return nil, err
    }
    
    // 检查缓存
    if result := ma.cache.Get(metric.Name()); result != nil {
        return result, nil
    }
    
    // 执行聚合计算
    result := &AggregateResult{
        Name:      metric.Name(),
        Timestamp: time.Now(),
    }
    
    switch rule.Operation {
    case SumOperation:
        result.Value = ma.calculateSum(window)
    case AvgOperation:
        result.Value = ma.calculateAvg(window)
    case MaxOperation:
        result.Value = ma.calculateMax(window)
    case MinOperation:
        result.Value = ma.calculateMin(window)
    case QuantileOperation:
        result.Quantiles = ma.calculateQuantiles(window, rule.Quantiles)
    }
    
    // 更新缓存
    ma.cache.Set(metric.Name(), result)
    
    return result, nil
}
