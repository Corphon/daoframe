// model/wuxing.go

package model

import (
    "sync"
    "time"
    "errors"
    
    "github.com/Corphon/daoframe/core"
)

var (
    ErrInvalidPhase = errors.New("无效的五行相位")
    ErrCycleBreak   = errors.New("五行循环中断")
)

// Phase 五行相位
type Phase uint8

const (
    PhaseWood Phase = iota // 木
    PhaseFire             // 火
    PhaseEarth            // 土
    PhaseMetal            // 金
    PhaseWater            // 水
)

// Relationship 五行关系
type Relationship uint8

const (
    RelGenerate Relationship = iota // 相生
    RelControl                      // 相克
    RelWeaken                       // 相泄
    RelNeutral                      // 中性
)

// Element 五行元素
type Element struct {
    mu         sync.RWMutex
    phase      Phase
    strength   uint8  // 0-100
    yinYang    *YinYang
    lastUpdate time.Time
}

// WuXing 五行系统
type WuXing struct {
    mu       sync.RWMutex
    elements map[Phase]*Element
    ctx      *core.DaoContext
    cycles   chan struct{}
    done     chan struct{}
}

// NewWuXing 创建新的五行系统
func NewWuXing(ctx *core.DaoContext) *WuXing {
    wx := &WuXing{
        elements: make(map[Phase]*Element),
        ctx:      ctx,
        cycles:   make(chan struct{}, 1),
        done:     make(chan struct{}),
    }

    // 初始化五行元素
    wx.initElements()
    
    // 启动五行循环
    go wx.runCycles()
    
    return wx
}

// initElements 初始化五行元素
func (wx *WuXing) initElements() {
    phases := []Phase{PhaseWood, PhaseFire, PhaseEarth, PhaseMetal, PhaseWater}
    
    for _, phase := range phases {
        wx.elements[phase] = &Element{
            phase:      phase,
            strength:   50, // 初始均衡
            yinYang:    NewYinYang(wx.ctx),
            lastUpdate: time.Now(),
        }
    }
}

// runCycles 运行五行循环
func (wx *WuXing) runCycles() {
    ticker := time.NewTicker(time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-wx.done:
            return
        case <-ticker.C:
            wx.processCycle()
        case <-wx.cycles:
            wx.processRelationships()
        }
    }
}

// processCycle 处理五行循环
func (wx *WuXing) processCycle() {
    wx.mu.Lock()
    defer wx.mu.Unlock()

    // 相生循环：木->火->土->金->水->木
    cycles := []Phase{PhaseWood, PhaseFire, PhaseEarth, PhaseMetal, PhaseWater}
    
    for i := 0; i < len(cycles); i++ {
        current := cycles[i]
        next := cycles[(i+1)%len(cycles)]
        
        // 促进相生关系
        wx.promote(current, next)
    }
}

// promote 促进两个元素间的相生关系
func (wx *WuXing) promote(from, to Phase) {
    source := wx.elements[from]
    target := wx.elements[to]

    source.mu.Lock()
    target.mu.Lock()
    defer source.mu.Unlock()
    defer target.mu.Unlock()

    // 根据源元素的强度增强目标元素
    if source.strength > 20 {
        energyTransfer := source.strength / 10
        source.strength -= energyTransfer
        target.strength += energyTransfer
        
        if target.strength > 100 {
            target.strength = 100
        }
    }
}

// GetRelationship 获取两个元素间的关系
func (wx *WuXing) GetRelationship(from, to Phase) Relationship {
    // 相生关系
    generates := map[Phase]Phase{
        PhaseWood:  PhaseFire,
        PhaseFire:  PhaseEarth,
        PhaseEarth: PhaseMetal,
        PhaseMetal: PhaseWater,
        PhaseWater: PhaseWood,
    }

    // 相克关系
    controls := map[Phase]Phase{
        PhaseWood:  PhaseEarth,
        PhaseEarth: PhaseWater,
        PhaseWater: PhaseFire,
        PhaseFire:  PhaseMetal,
        PhaseMetal: PhaseWood,
    }

    if generates[from] == to {
        return RelGenerate
    }
    if controls[from] == to {
        return RelControl
    }
    if generates[to] == from {
        return RelWeaken
    }
    return RelNeutral
}

// AdjustElement 调整元素强度
func (wx *WuXing) AdjustElement(phase Phase, delta int8) error {
    wx.mu.Lock()
    defer wx.mu.Unlock()

    element, exists := wx.elements[phase]
    if !exists {
        return ErrInvalidPhase
    }

    element.mu.Lock()
    defer element.mu.Unlock()

    newStrength := int16(element.strength) + int16(delta)
    if newStrength < 0 {
        newStrength = 0
    } else if newStrength > 100 {
        newStrength = 100
    }

    element.strength = uint8(newStrength)
    element.lastUpdate = time.Now()

    // 通知循环系统
    select {
    case wx.cycles <- struct{}{}:
    default:
    }

    return nil
}

// GetElementStrength 获取元素强度
func (wx *WuXing) GetElementStrength(phase Phase) (uint8, error) {
    wx.mu.RLock()
    defer wx.mu.RUnlock()

    element, exists := wx.elements[phase]
    if !exists {
        return 0, ErrInvalidPhase
    }

    element.mu.RLock()
    defer element.mu.RUnlock()
    return element.strength, nil
}

// Close 关闭五行系统
func (wx *WuXing) Close() error {
    close(wx.done)
    for _, element := range wx.elements {
        if element.yinYang != nil {
            element.yinYang.Close()
        }
    }
    return nil
}

// processRelationships 处理五行关系
func (wx *WuXing) processRelationships() {
    wx.mu.Lock()
    defer wx.mu.Unlock()

    // 处理相生和相克关系
    for phase, element := range wx.elements {
        // 获取所有与当前元素相关的关系
        for otherPhase, otherElement := range wx.elements {
            if phase == otherPhase {
                continue
            }

            relationship := wx.GetRelationship(phase, otherPhase)
            wx.applyRelationship(element, otherElement, relationship)
        }
    }
}

// applyRelationship 应用五行关系
func (wx *WuXing) applyRelationship(from, to *Element, rel Relationship) {
    from.mu.Lock()
    to.mu.Lock()
    defer from.mu.Unlock()
    defer to.mu.Unlock()

    switch rel {
    case RelGenerate:
        // 相生：增强目标元素
        energyTransfer := from.strength / 10
        if energyTransfer > 0 {
            from.strength -= energyTransfer / 2
            to.strength += energyTransfer
        }
    case RelControl:
        // 相克：削弱目标元素
        if from.strength > to.strength {
            diff := (from.strength - to.strength) / 5
            if diff > 0 {
                to.strength -= diff
            }
        }
    case RelWeaken:
        // 相泄：双方都减弱
        if from.strength > 0 {
            from.strength--
        }
        if to.strength > 0 {
            to.strength--
        }
    }

    // 确保范围有效
    if to.strength > 100 {
        to.strength = 100
    }
}
