//daoframe/core/origin.go

package core

import (
    "context"
    "github.com/Corphon/daoframe/core/state"  // 新的导入
    "github.com/Corphon/daoframe/core/force"  // 新的导入
)

// Origin 代表"道"的本源
type Origin struct {
    essence *DaoContext    // 道之精髓
    energy  *AdaptSystem   // 道之能量
    form    *BaseDaoSource // 道之形态
    components []Component
    mu         sync.RWMutex
    done       chan struct{}
}
type Component interface {
    Init(ctx context.Context) error
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
}
// 太极 - 表示最初的统一状态
type TaiJi struct {
    origin *Origin
    state  state.State
}

// 创建太极，实现"道生一"
func NewTaiJi() *TaiJi {
    origin := &Origin{
        essence: NewDaoContext(context.Background()),
        energy:  NewAdaptSystem(DefaultInterval),
        form:    NewBaseDaoSource(),
    }
    
    return &TaiJi{
        origin: origin,
        state:  state.StateInactive,
    }
}

// Generate 生成万物的起点，对应"道生一"
func (t *TaiJi) Generate() (*YinYang, error) {
    if t.state != state.StateInactive {
        return nil, ErrInvalidState
    }
    
    // 初始化本源
    if err := t.origin.form.Initialize(t.origin.essence); err != nil {
        return nil, err
    }
    
    // 激活能量系统
    if err := t.origin.energy.Start(t.origin.essence); err != nil {
        return nil, err
    }
    
    t.state = state.StateActive
    
    // 返回阴阳二气，为"一生二"做准备
    return NewYinYang(t.origin), nil
}

// 宇宙常数
const (
    DefaultInterval = time.Second * 1 // 基本时间单位
    MaximumForce   = 100             // 最大作用力
    MinimumForce   = 1               // 最小作用力
)
