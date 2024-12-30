//daoframe/core/origin.go

package core

import (
    "context"
)

// Origin 代表"道"的本源
type Origin struct {
    essence *DaoContext    // 道之精髓
    energy  *AdaptSystem   // 道之能量
    form    *BaseDaoSource // 道之形态
}

// 太极 - 表示最初的统一状态
type TaiJi struct {
    origin *Origin
    state  State
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
        state:  StateInactive,
    }
}

// Generate 生成万物的起点，对应"道生一"
func (t *TaiJi) Generate() (*YinYang, error) {
    if t.state != StateInactive {
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
    
    t.state = StateActive
    
    // 返回阴阳二气，为"一生二"做准备
    return NewYinYang(t.origin), nil
}

// 宇宙常数
const (
    DefaultInterval = time.Second * 1 // 基本时间单位
    MaximumForce   = 100             // 最大作用力
    MinimumForce   = 1               // 最小作用力
)

// Force 表示作用力
type Force uint8

// 定义基本作用力
const (
    ForceCreate Force = iota // 生之力
    ForceDestroy            // 灭之力
    ForceTransform          // 变之力
    ForceBalance           // 衡之力
)
