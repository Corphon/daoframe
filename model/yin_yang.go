// model/yin_yang.go

package model

import (
    "sync"
    "time"
    "errors"
    
    "github.com/Corphon/daoframe/core"
)

var (
    ErrImbalance = errors.New("阴阳失衡")
    ErrExtreme   = errors.New("阴阳极端")
)

// Nature 定义事物的性质
type Nature uint8

const (
    NatureYin  Nature = iota // 阴性
    NatureYang               // 阳性
    NatureTai                // 太极
)

// Polarity 定义阴阳极性
type Polarity struct {
    Value    uint8  // 0-100
    Nature   Nature
    LastSync time.Time
}

// YinYang 阴阳结构体
type YinYang struct {
    mu      sync.RWMutex
    yin     Polarity
    yang    Polarity
    ctx     *core.DaoContext
    state   core.State
    
    // 变化速率 (每秒)
    changeRate float64
    
    // 事件通道
    changes chan struct{}
    done    chan struct{}
}

// NewYinYang 创建新的阴阳实例
func NewYinYang(ctx *core.DaoContext) *YinYang {
    yy := &YinYang{
        ctx: ctx,
        yin: Polarity{
            Value:    50,
            Nature:   NatureYin,
            LastSync: time.Now(),
        },
        yang: Polarity{
            Value:    50,
            Nature:   NatureYang,
            LastSync: time.Now(),
        },
        state:      core.StateActive,
        changeRate: 1.0,
        changes:    make(chan struct{}, 1),
        done:       make(chan struct{}),
    }

    go yy.autoBalance()
    return yy
}

// autoBalance 自动平衡协程
func (yy *YinYang) autoBalance() {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-yy.done:
            return
        case <-ticker.C:
            yy.balance()
        case <-yy.changes:
            yy.notifyChange()
        }
    }
}

// balance 执行阴阳平衡
func (yy *YinYang) balance() {
    yy.mu.Lock()
    defer yy.mu.Unlock()

    now := time.Now()
    hour := now.Hour()

    // 昼夜自然变化
    if hour >= 6 && hour < 18 {
        // 白天：阳升阴降
        yy.adjustPolarity(&yy.yang, &yy.yin)
    } else {
        // 夜晚：阴升阳降
        yy.adjustPolarity(&yy.yin, &yy.yang)
    }
}

// adjustPolarity 调整阴阳极性
func (yy *YinYang) adjustPolarity(rise, fall *Polarity) {
    // 基于变化速率计算调整量
    delta := yy.changeRate

    if rise.Value < 100 {
        rise.Value++
    }
    if fall.Value > 0 {
        fall.Value--
    }

    rise.LastSync = time.Now()
    fall.LastSync = time.Now()
}

// GetRatio 获取阴阳比例
func (yy *YinYang) GetRatio() (yin, yang float64) {
    yy.mu.RLock()
    defer yy.mu.RUnlock()

    total := float64(yy.yin.Value + yy.yang.Value)
    if total == 0 {
        return 0.5, 0.5
    }
    return float64(yy.yin.Value) / total, float64(yy.yang.Value) / total
}

// Adjust 手动调整阴阳值
func (yy *YinYang) Adjust(yinDelta, yangDelta int8) error {
    yy.mu.Lock()
    defer yy.mu.Unlock()

    newYin := int16(yy.yin.Value) + int16(yinDelta)
    newYang := int16(yy.yang.Value) + int16(yangDelta)

    if newYin < 0 || newYin > 100 || newYang < 0 || newYang > 100 {
        return ErrExtreme
    }

    yy.yin.Value = uint8(newYin)
    yy.yang.Value = uint8(newYang)
    
    // 通知变化
    select {
    case yy.changes <- struct{}{}:
    default:
    }

    return nil
}

// Transform 阴阳转化
func (yy *YinYang) Transform() {
    yy.mu.Lock()
    defer yy.mu.Unlock()

    // 交换阴阳值
    yy.yin.Value, yy.yang.Value = yy.yang.Value, yy.yin.Value
}

// Split 阴阳分离，用于"二生三"
func (yy *YinYang) Split() (*YinYang, *YinYang) {
    yy.mu.RLock()
    defer yy.mu.RUnlock()

    // 创建偏阴实例
    yinCtx := yy.ctx.Clone()
    yinInstance := NewYinYang(yinCtx)
    yinInstance.yin.Value = 70
    yinInstance.yang.Value = 30

    // 创建偏阳实例
    yangCtx := yy.ctx.Clone()
    yangInstance := NewYinYang(yangCtx)
    yangInstance.yin.Value = 30
    yangInstance.yang.Value = 70

    return yinInstance, yangInstance
}

// Close 关闭并清理资源
func (yy *YinYang) Close() error {
    close(yy.done)
    return nil
}

// IsBalanced 检查阴阳是否平衡
func (yy *YinYang) IsBalanced() bool {
    yy.mu.RLock()
    defer yy.mu.RUnlock()

    diff := int(yy.yin.Value) - int(yy.yang.Value)
    if diff < 0 {
        diff = -diff
    }
    return diff <= 10
}

// GetDominant 获取主导属性
func (yy *YinYang) GetDominant() Nature {
    yy.mu.RLock()
    defer yy.mu.RUnlock()

    if yy.yin.Value > yy.yang.Value {
        return NatureYin
    } else if yy.yang.Value > yy.yin.Value {
        return NatureYang
    }
    return NatureTai
}

// notifyChange 通知变化
func (yy *YinYang) notifyChange() {
    // 这里可以添加观察者模式，通知其他组件阴阳变化
}
