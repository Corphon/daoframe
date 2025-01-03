//system/universe.go
package system

// Universe 宇宙系统
type Universe struct {
     mu           sync.RWMutex
    ctx          *core.DaoContext
    config       *SystemConfig
    
    // 核心系统组件
    timeSystem   *TimeSystem
    bagua        *model.BaGua
    wuXing       *model.WuXing
    yinYang      *model.YinYang
    lifecycle    *LifecycleManager
    
    // 监控和控制
    state        SystemState
    metrics      *UniverseMetrics
    eventBus     *EventBus
    scheduler    *Scheduler
    
    // 系统同步
    wg           sync.WaitGroup
    done         chan struct{}
}

// UniverseMetrics 宇宙系统指标
type UniverseMetrics struct {
    StartTime        time.Time
    CycleCount       uint64
    InteractionCount uint64
    ErrorCount       uint64
    LastError        error
    Components       map[string]*ComponentMetrics
}

// 新增方法
func (u *Universe) Start(ctx context.Context) error {
    u.mu.Lock()
    defer u.mu.Unlock()
    
    if u.state != SystemStateInactive {
        return ErrInvalidState
    }
    
    u.state = SystemStateStarting
    
    // 初始化组件
    if err := u.initializeComponents(); err != nil {
        return err
    }
    
    // 启动调度器
    if err := u.scheduler.Start(); err != nil {
        return err
    }
    
    u.state = SystemStateRunning
    return nil
}

// Evolution 演化控制
func (u *Universe) Evolution() error {
    // 1. 时空演化
    if err := u.timeSystem.Progress(); err != nil {
        return err
    }

    // 2. 八卦能量流动
    u.bagua.ProcessEnergyFlows()

    // 3. 五行相互作用
    u.wuXing.ProcessInteractions()

    // 4. 阴阳平衡调节
    u.yinYang.Balance()

    // 5. 生命周期更新
    return u.lifecycle.Update()
}

// InteractionSystem 交互系统
type InteractionSystem struct {
    baguaEffects      map[Trigram][]Effect
    elementInfluences map[Phase]map[Phase]float64
    temporalPatterns  map[GanZhiPair]Pattern
}
