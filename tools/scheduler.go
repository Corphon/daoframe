// tools/scheduler.go

package tools

import (
    "context"
    "sync"
    "time"
    "errors"
)

var (
    ErrTaskExists     = errors.New("任务已存在")
    ErrTaskNotFound   = errors.New("任务不存在")
    ErrSchedulerStopped = errors.New("调度器已停止")
)

// TaskFunc 任务函数类型
type TaskFunc func(ctx context.Context) error

// TaskPriority 任务优先级
type TaskPriority int

const (
    PriorityLow TaskPriority = iota
    PriorityNormal
    PriorityHigh
    PriorityCritical
)

// TaskStatus 任务状态
type TaskStatus int

const (
    StatusPending TaskStatus = iota
    StatusRunning
    StatusCompleted
    StatusFailed
    StatusCancelled
)

// Task 任务结构
type Task struct {
    ID          string
    Name        string
    Func        TaskFunc
    Priority    TaskPriority
    Status      TaskStatus
    Interval    time.Duration  // 定时任务间隔
    NextRun     time.Time
    LastRun     time.Time
    Error       error
    Context     context.Context
    Cancel      context.CancelFunc
}

// DaoScheduler 调度器
type DaoScheduler struct {
    mu          sync.RWMutex
    tasks       map[string]*Task
    taskQueue   chan *Task
    workerPool  chan struct{}
    maxWorkers  int
    ctx         context.Context
    cancel      context.CancelFunc
    wg          sync.WaitGroup
    running     bool
}

// NewDaoScheduler 创建新的调度器
func NewDaoScheduler(maxWorkers int) *DaoScheduler {
    ctx, cancel := context.WithCancel(context.Background())
    
    ds := &DaoScheduler{
        tasks:      make(map[string]*Task),
        taskQueue:  make(chan *Task, 100),
        workerPool: make(chan struct{}, maxWorkers),
        maxWorkers: maxWorkers,
        ctx:        ctx,
        cancel:     cancel,
        running:    false,
    }

    return ds
}

// Start 启动调度器
func (ds *DaoScheduler) Start() {
    ds.mu.Lock()
    if ds.running {
        ds.mu.Unlock()
        return
    }
    ds.running = true
    ds.mu.Unlock()

    // 启动任务分发器
    go ds.dispatcher()

    // 启动定时任务检查器
    go ds.timerChecker()
}

// dispatcher 任务分发器
func (ds *DaoScheduler) dispatcher() {
    for {
        select {
        case <-ds.ctx.Done():
            return
        case task := <-ds.taskQueue:
            // 获取worker槽位
            ds.workerPool <- struct{}{}
            ds.wg.Add(1)
            
            go func(t *Task) {
                defer func() {
                    <-ds.workerPool // 释放worker槽位
                    ds.wg.Done()
                }()
                
                ds.executeTask(t)
            }(task)
        }
    }
}

// timerChecker 定时任务检查器
func (ds *DaoScheduler) timerChecker() {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ds.ctx.Done():
            return
        case <-ticker.C:
            ds.checkScheduledTasks()
        }
    }
}

// checkScheduledTasks 检查定时任务
func (ds *DaoScheduler) checkScheduledTasks() {
    ds.mu.RLock()
    now := time.Now()
    for _, task := range ds.tasks {
        if task.Interval > 0 && now.After(task.NextRun) {
            // 提交任务到队列
            ds.submitTask(task)
            // 更新下次运行时间
            task.NextRun = now.Add(task.Interval)
        }
    }
    ds.mu.RUnlock()
}

// AddTask 添加任务
func (ds *DaoScheduler) AddTask(id string, name string, fn TaskFunc, priority TaskPriority) error {
    ds.mu.Lock()
    defer ds.mu.Unlock()

    if _, exists := ds.tasks[id]; exists {
        return ErrTaskExists
    }

    taskCtx, taskCancel := context.WithCancel(ds.ctx)
    task := &Task{
        ID:       id,
        Name:     name,
        Func:     fn,
        Priority: priority,
        Status:   StatusPending,
        Context:  taskCtx,
        Cancel:   taskCancel,
    }

    ds.tasks[id] = task
    return nil
}

// AddScheduledTask 添加定时任务
func (ds *DaoScheduler) AddScheduledTask(id string, name string, fn TaskFunc, interval time.Duration) error {
    ds.mu.Lock()
    defer ds.mu.Unlock()

    if _, exists := ds.tasks[id]; exists {
        return ErrTaskExists
    }

    taskCtx, taskCancel := context.WithCancel(ds.ctx)
    task := &Task{
        ID:       id,
        Name:     name,
        Func:     fn,
        Priority: PriorityNormal,
        Status:   StatusPending,
        Interval: interval,
        NextRun:  time.Now().Add(interval),
        Context:  taskCtx,
        Cancel:   taskCancel,
    }

    ds.tasks[id] = task
    return nil
}

// RemoveTask 移除任务
func (ds *DaoScheduler) RemoveTask(id string) error {
    ds.mu.Lock()
    defer ds.mu.Unlock()

    task, exists := ds.tasks[id]
    if !exists {
        return ErrTaskNotFound
    }

    task.Cancel()
    delete(ds.tasks, id)
    return nil
}

// executeTask 执行任务
func (ds *DaoScheduler) executeTask(task *Task) {
    ds.mu.Lock()
    task.Status = StatusRunning
    task.LastRun = time.Now()
    ds.mu.Unlock()

    err := task.Func(task.Context)

    ds.mu.Lock()
    if err != nil {
        task.Status = StatusFailed
        task.Error = err
    } else {
        task.Status = StatusCompleted
        task.Error = nil
    }
    ds.mu.Unlock()
}

// submitTask 提交任务到队列
func (ds *DaoScheduler) submitTask(task *Task) {
    select {
    case ds.taskQueue <- task:
        // 任务成功提交到队列
    default:
        // 队列已满，记录错误
        DefaultLogger.Error("Task queue is full, task %s dropped", task.ID)
    }
}

// GetTaskStatus 获取任务状态
func (ds *DaoScheduler) GetTaskStatus(id string) (TaskStatus, error) {
    ds.mu.RLock()
    defer ds.mu.RUnlock()

    task, exists := ds.tasks[id]
    if !exists {
        return StatusPending, ErrTaskNotFound
    }

    return task.Status, nil
}

// Stop 停止调度器
func (ds *DaoScheduler) Stop() {
    ds.mu.Lock()
    if !ds.running {
        ds.mu.Unlock()
        return
    }
    ds.running = false
    ds.cancel()
    ds.mu.Unlock()

    // 等待所有任务完成
    ds.wg.Wait()
}

// IsRunning 检查调度器是否运行中
func (ds *DaoScheduler) IsRunning() bool {
    ds.mu.RLock()
    defer ds.mu.RUnlock()
    return ds.running
}
