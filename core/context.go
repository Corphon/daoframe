// core/context.go

package core

import (
    "context"
    "sync"
    "time"
)

// DaoPhase 代表事物所处的阶段
type DaoPhase uint8

const (
    PhaseWuJi     DaoPhase = iota // 无极
    PhaseTaiJi                    // 太极
    PhaseYinYang                  // 阴阳
    PhaseBaGua                    // 八卦
    PhaseWanWu                    // 万物
)

// DaoAttribute 描述事物的属性
type DaoAttribute struct {
    Yin  uint8 // 阴性程度 (0-100)
    Yang uint8 // 阳性程度 (0-100)
}

// DaoContext 扩展标准 context，体现道的特性
type DaoContext struct {
    context.Context
    mu         sync.RWMutex
    phase      DaoPhase                 // 当前阶段
    attributes *DaoAttribute            // 阴阳属性
    values     map[string]interface{}   // 存储值
    birth      time.Time                // 创建时间
}

// NewDaoContext 创建新的道家上下文
func NewDaoContext(ctx context.Context) *DaoContext {
    if ctx == nil {
        ctx = context.Background()
    }
    
    return &DaoContext{
        Context:    ctx,
        phase:      PhaseWuJi,
        attributes: &DaoAttribute{Yin: 50, Yang: 50}, // 初始平衡
        values:     make(map[string]interface{}),
        birth:      time.Now(),
    }
}

// SetPhase 设置阶段
func (dc *DaoContext) SetPhase(p DaoPhase) {
    dc.mu.Lock()
    defer dc.mu.Unlock()
    dc.phase = p
}

// GetPhase 获取当前阶段
func (dc *DaoContext) GetPhase() DaoPhase {
    dc.mu.RLock()
    defer dc.mu.RUnlock()
    return dc.phase
}

// SetValue 设置值
func (dc *DaoContext) SetValue(key string, value interface{}) {
    dc.mu.Lock()
    defer dc.mu.Unlock()
    dc.values[key] = value
}

// GetValue 获取值
func (dc *DaoContext) GetValue(key string) (interface{}, bool) {
    dc.mu.RLock()
    defer dc.mu.RUnlock()
    value, exists := dc.values[key]
    return value, exists
}

// AdjustAttribute 调整阴阳属性
func (dc *DaoContext) AdjustAttribute(yinDelta, yangDelta int16) {
    dc.mu.Lock()
    defer dc.mu.Unlock()
    
    // 确保阴阳值在 0-100 范围内
    newYin := int16(dc.attributes.Yin) + yinDelta
    newYang := int16(dc.attributes.Yang) + yangDelta
    
    if newYin < 0 {
        newYin = 0
    } else if newYin > 100 {
        newYin = 100
    }
    
    if newYang < 0 {
        newYang = 0
    } else if newYang > 100 {
        newYang = 100
    }

    // 添加阴阳平衡检查
    if !isBalanceValid(uint8(newYin), uint8(newYang)) {
        return ErrInvalidBalance
    }
    
    dc.attributes.Yin = uint8(newYin)
    dc.attributes.Yang = uint8(newYang)
}

// [新增] 阴阳平衡检查
func isBalanceValid(yin, yang uint8) bool {
    diff := int(yin) - int(yang)
    if diff < 0 {
        diff = -diff
    }
    return diff <= 30 // 允许适度的不平衡
}

// GetAttribute 获取阴阳属性
func (dc *DaoContext) GetAttribute() DaoAttribute {
    dc.mu.RLock()
    defer dc.mu.RUnlock()
    return *dc.attributes
}

// Age 获取上下文年龄
func (dc *DaoContext) Age() time.Duration {
    return time.Since(dc.birth)
}

// WithTimeout 创建具有超时的新上下文
func (dc *DaoContext) WithTimeout(timeout time.Duration) (*DaoContext, context.CancelFunc) {
    ctx, cancel := context.WithTimeout(dc.Context, timeout)
    newCtx := NewDaoContext(ctx)
    // 继承原上下文的属性
    newCtx.phase = dc.phase
    newCtx.attributes = &DaoAttribute{
        Yin:  dc.attributes.Yin,
        Yang: dc.attributes.Yang,
    }
    return newCtx, cancel
}

// WithCancel 创建可取消的新上下文
func (dc *DaoContext) WithCancel() (*DaoContext, context.CancelFunc) {
    ctx, cancel := context.WithCancel(dc.Context)
    newCtx := NewDaoContext(ctx)
    // 继承原上下文的属性
    newCtx.phase = dc.phase
    newCtx.attributes = &DaoAttribute{
        Yin:  dc.attributes.Yin,
        Yang: dc.attributes.Yang,
    }
    return newCtx, cancel
}

// Clone 创建上下文的克隆
func (dc *DaoContext) Clone() *DaoContext {
    dc.mu.RLock()
    defer dc.mu.RUnlock()
    
    newCtx := NewDaoContext(dc.Context)
    newCtx.phase = dc.phase
    newCtx.attributes = &DaoAttribute{
        Yin:  dc.attributes.Yin,
        Yang: dc.attributes.Yang,
    }
    
    // 复制值
    for k, v := range dc.values {
        newCtx.values[k] = v
    }
    
    return newCtx
}
