// tools/scheduler/scheduler.go
package scheduler

// Scheduler 调度器
type Scheduler struct {
    tasks       map[string]*Task
    pool        *async.WorkerPool
    timer       *time.Timer
    dispatcher  *TaskDispatcher
    manager     *TaskManager
    metrics     *SchedulerMetrics
    mu          sync.RWMutex
}

// TaskDispatcher 任务分发器
type TaskDispatcher struct {
    queues     map[TaskPriority]*TaskQueue
    selector   TaskSelector
    limiter    rate.Limiter
    metrics    *DispatcherMetrics
}

// TaskManager 任务管理器
type TaskManager struct {
    scheduler  *Scheduler
    store      TaskStore
    hooks      []TaskHook
    recovery   TaskRecovery
}

// TaskStore 任务存储接口
type TaskStore interface {
    Save(task *Task) error
    Load(id string) (*Task, error)
    List() ([]*Task, error)
    Delete(id string) error
}
