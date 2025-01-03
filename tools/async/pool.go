// tools/async/pool.go
package async

// WorkerPool 工作池
type WorkerPool struct {
    workers    []*Worker
    taskQueue  chan Task
    results    chan Result
    workerSize int
    mu         sync.RWMutex
    wg         sync.WaitGroup
    ctx        context.Context
    cancel     context.CancelFunc
}

// Worker 工作者
type Worker struct {
    id      int
    pool    *WorkerPool
    metrics *WorkerMetrics
}

// Task 任务接口
type Task interface {
    Execute(context.Context) (interface{}, error)
    ID() string
    Priority() int
}

// WorkerMetrics 工作者指标
type WorkerMetrics struct {
    TasksProcessed uint64
    ErrorCount     uint64
    ProcessingTime time.Duration
    IdleTime       time.Duration
}
