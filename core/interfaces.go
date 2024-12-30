// core/interfaces.go

package core

import (
    "context"
    "time"
)

// DaoEntity 代表一个道的实体
type DaoEntity interface {
    // GetID 获取实体标识
    GetID() string
    
    // GetState 获取当前状态
    GetState() State
    
    // GetAttributes 获取实体属性
    GetAttributes() map[string]interface{}
}

// Force 表示作用力的类型
type Force uint8

const (
    ForceNone     Force = iota // 无力
    ForceYin                   // 阴力
    ForceYang                  // 阳力
    ForceBalance              // 平衡力
    ForceChange               // 变化力
)

// State 表示状态
type State uint8

const (
    StateVoid     State = iota // 虚无态
    StateInactive              // 静止态
    StateActive               // 活跃态
    StatePaused               // 停滞态
    StateChanged              // 变化态
    StateTerminated          // 终止态
)

// Phase 表示五行相位
type Phase uint8

const (
    PhaseWood  Phase = iota // 木
    PhaseFire               // 火
    PhaseEarth              // 土
    PhaseMetal              // 金
    PhaseWater              // 水
)

// DaoSource 定义道源接口
type DaoSource interface {
    // Initialize 初始化，从虚无中生成
    Initialize(ctx context.Context) error
    
    // Activate 激活，使之具有生机
    Activate(ctx context.Context) error
    
    // Adapt 适应环境变化，体现道的自然特性
    Adapt(ctx context.Context) error
    
    // ApplyForce 施加作用力，引发变化
    ApplyForce(force Force) error
    
    // GetState 获取当前状态
    GetState() State
    
    // Terminate 返归虚无
    Terminate(ctx context.Context) error
}

// Transformer 定义事物转化接口
type Transformer interface {
    // Transform 进行状态转换
    Transform(ctx context.Context, from, to State) error
    
    // CanTransform 检查是否可以转换
    CanTransform(from, to State) bool
    
    // GetTransformPath 获取转换路径
    GetTransformPath(from, to State) []State
}

// Observer 定义观察者接口
type Observer interface {
    // OnStateChange 状态变更通知
    OnStateChange(entity DaoEntity, oldState, newState State)
    
    // OnForceApplied 力作用通知
    OnForceApplied(entity DaoEntity, force Force)
}

// LifeCycleManager 定义生命周期管理接口
type LifeCycleManager interface {
    // CreateEntity 创建实体
    CreateEntity(ctx context.Context, id string) (DaoEntity, error)
    
    // DestroyEntity 销毁实体
    DestroyEntity(ctx context.Context, id string) error
    
    // GetEntity 获取实体
    GetEntity(id string) (DaoEntity, error)
    
    // RegisterObserver 注册观察者
    RegisterObserver(observer Observer)
    
    // RemoveObserver 移除观察者
    RemoveObserver(observer Observer)
}

// Scheduler 定义调度器接口
type Scheduler interface {
    // Schedule 调度任务
    Schedule(ctx context.Context, task Task) error
    
    // Cancel 取消任务
    Cancel(taskID string) error
    
    // GetTaskStatus 获取任务状态
    GetTaskStatus(taskID string) (TaskStatus, error)
}

// Task 定义任务接口
type Task interface {
    // GetID 获取任务ID
    GetID() string
    
    // Execute 执行任务
    Execute(ctx context.Context) error
    
    // GetPriority 获取优先级
    GetPriority() int
    
    // GetDeadline 获取截止时间
    GetDeadline() time.Time
}

// TaskStatus 定义任务状态
type TaskStatus struct {
    ID        string
    State     State
    Progress  float64
    Error     error
    StartTime time.Time
    EndTime   time.Time
}

// StateCoordinator 定义状态协调器接口
type StateCoordinator interface {
    // RegisterState 注册状态
    RegisterState(state State) error
    
    // AddTransition 添加状态转换规则
    AddTransition(from, to State) error
    
    // ValidateTransition 验证状态转换
    ValidateTransition(from, to State) bool
    
    // GetValidTransitions 获取有效的转换状态
    GetValidTransitions(state State) []State
}

// MetricsCollector 定义指标收集接口
type MetricsCollector interface {
    // RecordStateChange 记录状态变更
    RecordStateChange(entity DaoEntity, oldState, newState State)
    
    // RecordForceApplication 记录力的作用
    RecordForceApplication(entity DaoEntity, force Force)
    
    // RecordTaskExecution 记录任务执行
    RecordTaskExecution(task Task, duration time.Duration, err error)
    
    // GetMetrics 获取指标数据
    GetMetrics() map[string]interface{}
}
