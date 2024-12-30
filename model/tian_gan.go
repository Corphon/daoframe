// model/tian_gan.go

package model

import (
    "sync"
    "time"
    "errors"
    
    "github.com/Corphon/daoframe/core"
)

var (
    ErrInvalidGan = errors.New("无效的天干")
    ErrCycleFail  = errors.New("天干运转失败")
)

// Gan 天干类型
type Gan uint8

const (
    GanJia  Gan = iota // 甲
    GanYi              // 乙
    GanBing            // 丙
    GanDing            // 丁
    GanWu              // 戊
    GanJi              // 己
    GanGeng            // 庚
    GanXin             // 辛
    GanRen             // 壬
    GanGui             // 癸
)

// GanAttribute 天干属性
type GanAttribute struct {
    Element  Phase     // 所属五行
    Nature   Nature    // 阴阳属性
    Position uint8     // 位置(0-9)
    Energy   uint8     // 能量值(0-100)
}

// TianGan 天干系统
type TianGan struct {
    mu        sync.RWMutex
    gans      map[Gan]*GanAttribute
    current   Gan
    wuxing    *WuXing
    ctx       *core.DaoContext
    changes   chan Gan
    done      chan struct{}
}

// NewTianGan 创建天干系统
func NewTianGan(ctx *core.DaoContext, wx *WuXing) *TianGan {
    tg := &TianGan{
        gans:    make(map[Gan]*GanAttribute),
        wuxing:  wx,
        ctx:     ctx,
        changes: make(chan Gan, 1),
        done:    make(chan struct{}),
    }

    tg.initGans()
    go tg.run()
    return tg
}

// initGans 初始化天干
func (tg *TianGan) initGans() {
    // 天干五行属性对应
    elementMap := map[Gan]Phase{
        GanJia:  PhaseWood,  // 甲木
        GanYi:   PhaseWood,  // 乙木
        GanBing: PhaseFire,  // 丙火
        GanDing: PhaseFire,  // 丁火
        GanWu:   PhaseEarth, // 戊土
        GanJi:   PhaseEarth, // 己土
        GanGeng: PhaseMetal, // 庚金
        GanXin:  PhaseMetal, // 辛金
        GanRen:  PhaseWater, // 壬水
        GanGui:  PhaseWater, // 癸水
    }

    // 天干阴阳属性
    natureMap := map[Gan]Nature{
        GanJia:  NatureYang, // 阳
        GanYi:   NatureYin,  // 阴
        GanBing: NatureYang, // 阳
        GanDing: NatureYin,  // 阴
        GanWu:   NatureYang, // 阳
        GanJi:   NatureYin,  // 阴
        GanGeng: NatureYang, // 阳
        GanXin:  NatureYin,  // 阴
        GanRen:  NatureYang, // 阳
        GanGui:  NatureYin,  // 阴
    }

    for i := GanJia; i <= GanGui; i++ {
        tg.gans[i] = &GanAttribute{
            Element:  elementMap[i],
            Nature:   natureMap[i],
            Position: uint8(i),
            Energy:   50, // 初始能量均衡
        }
    }

    tg.current = GanJia // 从甲开始
}

// run 运行天干循环
func (tg *TianGan) run() {
    ticker := time.NewTicker(time.Hour * 2) // 每两小时转动一次
    defer ticker.Stop()

    for {
        select {
        case <-tg.done:
            return
        case <-ticker.C:
            tg.rotate()
        case gan := <-tg.changes:
            tg.handleChange(gan)
        }
    }
}

// rotate 天干轮转
func (tg *TianGan) rotate() {
    tg.mu.Lock()
    defer tg.mu.Unlock()

    // 计算下一个天干
    next := (tg.current + 1) % 10
    
    // 更新能量
    prevAttr := tg.gans[tg.current]
    nextAttr := tg.gans[next]
    
    // 能量传递
    energyTransfer := prevAttr.Energy / 10
    prevAttr.Energy -= energyTransfer
    nextAttr.Energy += energyTransfer

    // 更新当前天干
    tg.current = next
}

// GetCurrentGan 获取当前天干
func (tg *TianGan) GetCurrentGan() Gan {
    tg.mu.RLock()
    defer tg.mu.RUnlock()
    return tg.current
}

// GetGanAttribute 获取天干属性
func (tg *TianGan) GetGanAttribute(gan Gan) (*GanAttribute, error) {
    tg.mu.RLock()
    defer tg.mu.RUnlock()

    attr, exists := tg.gans[gan]
    if !exists {
        return nil, ErrInvalidGan
    }
    
    // 返回属性副本
    return &GanAttribute{
        Element:  attr.Element,
        Nature:   attr.Nature,
        Position: attr.Position,
        Energy:   attr.Energy,
    }, nil
}

// AdjustEnergy 调整天干能量
func (tg *TianGan) AdjustEnergy(gan Gan, delta int8) error {
    tg.mu.Lock()
    defer tg.mu.Unlock()

    attr, exists := tg.gans[gan]
    if !exists {
        return ErrInvalidGan
    }

    newEnergy := int16(attr.Energy) + int16(delta)
    if newEnergy < 0 {
        newEnergy = 0
    } else if newEnergy > 100 {
        newEnergy = 100
    }

    attr.Energy = uint8(newEnergy)

    // 通知变化
    select {
    case tg.changes <- gan:
    default:
    }

    return nil
}

// handleChange 处理天干变化
func (tg *TianGan) handleChange(gan Gan) {
    tg.mu.Lock()
    defer tg.mu.Unlock()

    attr := tg.gans[gan]
    
    // 影响对应的五行
    if tg.wuxing != nil {
        tg.wuxing.AdjustElement(attr.Element, int8(attr.Energy/10))
    }
}

// GetGanElement 获取天干对应的五行
func (tg *TianGan) GetGanElement(gan Gan) Phase {
    tg.mu.RLock()
    defer tg.mu.RUnlock()

    if attr, exists := tg.gans[gan]; exists {
        return attr.Element
    }
    return PhaseWood // 默认返回木
}

// GetGanNature 获取天干阴阳属性
func (tg *TianGan) GetGanNature(gan Gan) Nature {
    tg.mu.RLock()
    defer tg.mu.RUnlock()

    if attr, exists := tg.gans[gan]; exists {
        return attr.Nature
    }
    return NatureYang // 默认返回阳
}

// Close 关闭天干系统
func (tg *TianGan) Close() error {
    close(tg.done)
    return nil
}
