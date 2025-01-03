// core/adapt.go

package core

import (
    "context"
    "errors"
    "sync"
    "time"
    "github.com/Corphon/daoframe/core/state"  // 新的导入
)

// AdaptHandler 定义适应性处理函数类型
type AdaptHandler func(ctx *DaoContext) error

// AdaptMode 适应模式
type AdaptMode uint8

const (
    // 自然适应：渐进、温和
    NaturalAdapt AdaptMode = iota
    // 主动适应：积极、快速
    ActiveAdapt
    // 被动适应：保守、稳定
    PassiveAdapt
)

// AdaptSystem 实现自适应系统
type AdaptSystem struct {
    mu          sync.RWMutex
    stateManager *state.StateManager  // 使用状态管理器
    handlers    map[string]AdaptHandler
    mode        AdaptMode
    active      bool
    interval    time.Duration
    lastAdapt   time.Time
    yinHandler  []AdaptHandler
    yangHandler []AdaptHandler
    // 新增字段
    adaptiveThreshold float64
    environmentState  map[string]int
    adaptHistory     []AdaptiveAction
    balanceFactors   map[string]float64
    currentState state.State  // 新增字段用于跟踪系统状态
}

// NewAdaptSystem 创建新的自适应系统
func NewAdaptSystem(interval time.Duration) *AdaptSystem {
    return &AdaptSystem{
        handlers:    make(map[string]AdaptHandler),
        mode:       NaturalAdapt,
        interval:   interval,
        yinHandler: make([]AdaptHandler, 0),
        yangHandler: make([]AdaptHandler, 0),
        currentState: state.StateInactive,  // 初始化状态
    }
}

// RegisterHandler 注册处理器
func (as *AdaptSystem) RegisterHandler(name string, handler AdaptHandler, nature DaoPhase) error {
    if handler == nil {
        return errors.New("handler cannot be nil")
    }

    as.mu.Lock()
    defer as.mu.Unlock()

    as.handlers[name] = handler
    
    // 根据性质分类处理器
    switch nature {
    case PhaseYinYang:
        if as.isYinDominant() {
            as.yinHandler = append(as.yinHandler, handler)
        } else {
            as.yangHandler = append(as.yangHandler, handler)
        }
    default:
        // 其他阶段的处理器保持中性
    }

    return nil
}

// SetMode 设置适应模式
func (as *AdaptSystem) SetMode(mode AdaptMode) {
    as.mu.Lock()
    defer as.mu.Unlock()
    as.mode = mode
}

// Start 启动自适应系统
func (as *AdaptSystem) Start(ctx context.Context) error {
    as.mu.Lock()
    defer as.mu.Unlock()

    if err := as.stateManager.TransitTo(state.StateActive); err != nil {
        return fmt.Errorf("failed to start adapt system: %w", err)
    }
    if as.currentState == state.StateActive {
        as.mu.Unlock()
        return errors.New("adapt system is already running")
    }
    
    if !state.ValidateTransition(as.currentState, state.StateActive) {
        as.mu.Unlock()
        return errors.New("invalid state transition")
    }
    
    as.active = true
    as.currentState = state.StateActive
    as.lastAdapt = time.Now()
    as.mu.Unlock()

    go as.run(ctx)
    return nil
}
// Stop 停止自适应系统
func (as *AdaptSystem) Stop() {
    as.mu.Lock()
    defer as.mu.Unlock()
    
    if state.ValidateTransition(as.currentState, state.StateInactive) {
        as.active = false
        as.currentState = state.StateInactive
    }
}

// run 运行自适应循环
func (as *AdaptSystem) run(ctx context.Context) {
    ticker := time.NewTicker(as.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            as.Stop()
            return
        case <-ticker.C:
            if err := as.adapt(ctx); err != nil {
                continue
            }
        }
    }
}

// adapt 执行适应过程
func (as *AdaptSystem) adapt(ctx context.Context) error {
    as.mu.RLock()
    mode := as.mode
    handlers := make([]struct {
        name    string
        handler AdaptHandler
    }, 0, len(as.handlers))

    for name, handler := range as.handlers {
        handlers = append(handlers, struct {
            name    string
            handler AdaptHandler
        }{name, handler})
    }
    as.mu.RUnlock()

    daoCtx := NewDaoContext(ctx)

    switch mode {
    case NaturalAdapt:
        return as.naturalAdapt(daoCtx, handlers)
    case ActiveAdapt:
        return as.activeAdapt(daoCtx, handlers)
    default:
        return as.passiveAdapt(daoCtx, handlers)
    }
}

// naturalAdapt 自然适应过程
func (as *AdaptSystem) naturalAdapt(ctx *DaoContext, handlers []struct {
    name    string
    handler AdaptHandler
}) error {
    // 遵循自然规律，平衡阴阳
    for _, h := range handlers {
        if err := h.handler(ctx); err != nil {
            continue
        }
        // 自然间隔
        time.Sleep(time.Millisecond * 100)
    }
    return nil
}

// activeAdapt 主动适应过程
func (as *AdaptSystem) activeAdapt(ctx *DaoContext, handlers []struct {
    name    string
    handler AdaptHandler
}) error {
    // 并发执行，快速适应
    var wg sync.WaitGroup
    for _, h := range handlers {
        wg.Add(1)
        go func(handler AdaptHandler) {
            defer wg.Done()
            _ = handler(ctx)
        }(h.handler)
    }
    wg.Wait()
    return nil
}

// passiveAdapt 被动适应过程
func (as *AdaptSystem) passiveAdapt(ctx *DaoContext, handlers []struct {
    name    string
    handler AdaptHandler
}) error {
    // 保守执行，注重稳定
    for _, h := range handlers {
        if err := h.handler(ctx); err != nil {
            return err // 遇错即停
        }
        // 较长间隔
        time.Sleep(time.Millisecond * 200)
    }
    return nil
}

// isYinDominant 判断是否阴性主导
func (as *AdaptSystem) isYinDominant() bool {
    return time.Now().Hour() >= 18 || time.Now().Hour() < 6
}

// getAdaptInterval 获取适应间隔
func (as *AdaptSystem) getAdaptInterval() time.Duration {
    switch as.mode {
    case ActiveAdapt:
        return as.interval / 2
    case PassiveAdapt:
        return as.interval * 2
    default:
        return as.interval
    }
}
