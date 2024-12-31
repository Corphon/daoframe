// core/dao.go

package core

import (
    "context"
    "errors"
)

// 定义核心错误类型
var (
    ErrInvalidState   = errors.New("invalid state")
    ErrInitFailed    = errors.New("initialization failed")
    ErrNotInitialized = errors.New("dao source not initialized")
    ErrInvalidForce   = errors.New("invalid force application")
)

// Force 表示道的作用力
type Force uint8

const (
    ForceCreate  Force = iota + 1 // 生之力
    ForceDestroy                  // 灭之力
    ForceTransform                // 变之力
    ForceBalance                  // 衡之力
)

// DaoSource 定义了道源的核心接口
type DaoSource interface {
    // Initialize 初始化道源，从虚无中生一
    Initialize(ctx context.Context) error
    
    // Adapt 适应环境变化，体现道的自然特性
    Adapt(ctx context.Context) error
    
    // ApplyForce 施加作用力，引发变化
    ApplyForce(force Force) error
    
    // Terminate 返归虚无
    Terminate(ctx context.Context) error
    
    // GetState 获取当前状态
    GetState() State
}

// BaseDaoSource 提供 DaoSource 接口的基本实现
type BaseDaoSource struct {
    state     State
    essence   interface{} // 本质
    force     Force      // 当前作用力
}

// NewBaseDaoSource 创建新的道源基础实现
func NewBaseDaoSource() *BaseDaoSource {
    return &BaseDaoSource{
        state: StateVoid,
        force: ForceCreate,
    }
}

// Initialize 实现基本的初始化
func (b *BaseDaoSource) Initialize(ctx context.Context) error {
    if b.state != StateVoid {
        return ErrInvalidState
    }
    b.state = StateInactive
    return nil
}

// Activate 激活道源
func (b *BaseDaoSource) Activate(ctx context.Context) error {
    if b.state != StateInactive {
        return ErrInvalidState
    }
    b.state = StateActive
    return nil
}

// ApplyForce 实现力的作用
func (b *BaseDaoSource) ApplyForce(force Force) error {
    if b.state != StateActive {
        return ErrInvalidState
    }
    
    switch force {
    case ForceCreate:
        if b.force == ForceDestroy {
            return ErrInvalidForce
        }
    case ForceDestroy:
        if b.force == ForceCreate {
            return ErrInvalidForce
        }
    case ForceTransform:
        // 变化之力可以在任何状态下使用
    case ForceBalance:
        // 平衡之力可以化解对立
        b.force = ForceBalance
    default:
        return ErrInvalidForce
    }
    
    b.force = force
    return nil
}

// Adapt 实现基本的适应机制
func (b *BaseDaoSource) Adapt(ctx context.Context) error {
    if b.state != StateActive {
        return ErrInvalidState
    }
    return nil
}

// Terminate 实现基本的终止逻辑
func (b *BaseDaoSource) Terminate(ctx context.Context) error {
    if b.state == StateTerminated {
        return ErrInvalidState
    }
    b.state = StateTerminated
    return nil
}

// GetState 获取当前状态
func (b *BaseDaoSource) GetState() State {
    return b.state
}

// GetForce 获取当前作用力
func (b *BaseDaoSource) GetForce() Force {
    return b.force
}

// isValidTransition 检查状态转换是否有效
func isValidTransition(current, new State) bool {
    switch current {
    case StateVoid:
        return new == StateInactive
    case StateInactive:
        return new == StateActive || new == StateVoid
    case StateActive:
        return new == StatePaused || new == StateTerminated || new == StateInactive
    case StatePaused:
        return new == StateActive || new == StateTerminated
    case StateTerminated:
        return new == StateVoid // 允许重新归于虚无
    default:
        return false
    }
}
