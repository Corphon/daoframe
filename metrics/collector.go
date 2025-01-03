//metrics/collector.go
package metrics

import (
    "context"
    "sync"
    "time"
    
    "github.com/Corphon/daoframe/errors"
    "github.com/Corphon/daoframe/tools/async"
)

// MetricCollector 指标收集器
type MetricCollector struct {
    mu          sync.RWMutex
    metrics     map[string]Metric
    labels      map[string]string
    aggregator  *MetricAggregator
    reporter    *MetricReporter
    storage     MetricStorage
    buffer      *ring.Buffer
    
    ctx         context.Context
    cancel      context.CancelFunc
    wg          sync.WaitGroup
    
    // 配置选项
    opts        CollectorOptions
    // 工作池
    pool        *async.WorkerPool
    // 采样控制
    sampler     *MetricSampler
    // 监控指标
    metrics     *CollectorMetrics
}

// CollectorOptions 收集器配置
type CollectorOptions struct {
    // 缓冲区大小
    BufferSize int
    // 工作池大小
    WorkerSize int
    // 采样率
    SampleRate float64
    // 刷新间隔
    FlushInterval time.Duration
    // 存储选项
    StorageOptions StorageOptions
    // 聚合选项
    AggregateOptions AggregateOptions
    // 报告选项
    ReportOptions ReportOptions
}

// NewMetricCollector 创建指标收集器
func NewMetricCollector(opts CollectorOptions) (*MetricCollector, error) {
    if err := validateOptions(opts); err != nil {
        return nil, err
    }
    
    ctx, cancel := context.WithCancel(context.Background())
    
    mc := &MetricCollector{
        metrics:    make(map[string]Metric),
        labels:     make(map[string]string),
        buffer:     ring.New(opts.BufferSize),
        ctx:        ctx,
        cancel:     cancel,
        opts:       opts,
    }
    
    // 初始化工作池
    pool, err := async.NewWorkerPool(async.WorkerPoolOptions{
        Size:      opts.WorkerSize,
        QueueSize: opts.BufferSize,
    })
    if err != nil {
        cancel()
        return nil, err
    }
    mc.pool = pool
    
    // 初始化采样器
    mc.sampler = NewMetricSampler(opts.SampleRate)
    
    // 初始化存储
    storage, err := NewMetricStorage(opts.StorageOptions)
    if err != nil {
        cancel()
        return nil, err
    }
    mc.storage = storage
    
    // 初始化聚合器
    mc.aggregator = NewMetricAggregator(opts.AggregateOptions)
    
    // 初始化报告器
    mc.reporter = NewMetricReporter(opts.ReportOptions)
    
    // 启动后台任务
    mc.startTasks()
    
    return mc, nil
}

// RegisterMetric 注册指标
func (mc *MetricCollector) RegisterMetric(name string, typ MetricType, labels map[string]string) (Metric, error) {
    mc.mu.Lock()
    defer mc.mu.Unlock()
    
    // 检查是否已存在
    if _, exists := mc.metrics[name]; exists {
        return nil, errors.New("metric already exists")
    }
    
    // 创建新指标
    metric, err := mc.createMetric(name, typ, labels)
    if err != nil {
        return nil, err
    }
    
    // 注册指标
    mc.metrics[name] = metric
    
    // 通知聚合器
    mc.aggregator.OnMetricRegistered(metric)
    
    return metric, nil
}

// Collect 收集指标值
func (mc *MetricCollector) Collect(name string, value float64, labels map[string]string) error {
    // 采样控制
    if !mc.sampler.ShouldSample() {
        return nil
    }
    
    // 获取指标
    metric, err := mc.getMetric(name)
    if err != nil {
        return err
    }
    
    // 异步处理
    return mc.pool.Submit(func() error {
        // 更新指标值
        metric.Add(value)
        
        // 写入缓冲区
        mc.buffer.Value = &MetricPoint{
            Name:   name,
            Value:  value,
            Labels: labels,
            Time:   time.Now(),
        }
        mc.buffer = mc.buffer.Next()
        
        return nil
    })
}
