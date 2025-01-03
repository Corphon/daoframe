//system/scheduler.go
package system

type Scheduler struct {
    mu        sync.RWMutex
    tasks     map[string]*Task
    workers   []*Worker
    queue     chan *Task
    metrics   *SchedulerMetrics
    
    config    *SchedulerConfig
    state     SystemState
    done      chan struct{}
}

type Task struct {
    ID          string
    Priority    int
    Interval    time.Duration
    Action      func(context.Context) error
    LastRun     time.Time
    NextRun     time.Time
    Stats       *TaskStats
}

type SchedulerConfig struct {
    WorkerCount     int
    QueueSize       int
    MaxRetries      int
    RetryDelay      time.Duration
    DefaultPriority int
}

func NewScheduler(config *SchedulerConfig) *Scheduler {
    return &Scheduler{
        tasks:   make(map[string]*Task),
        queue:   make(chan *Task, config.QueueSize),
        workers: make([]*Worker, config.WorkerCount),
        config:  config,
        done:    make(chan struct{}),
    }
}
