// model/di_zhi.go

package model

import (
    "sync"
    "time"
    "errors"
    
    "github.com/Corphon/daoframe/core"
)

var (
    ErrInvalidZhi    = errors.New("无效的地支")
    ErrZhiLocked     = errors.New("地支已锁定")
    ErrInvalidCycle  = errors.New("无效的周期")
)

// Zhi 地支枚举
type Zhi int

const (
    ZhiZi   Zhi = iota // 子
    ZhiChou            // 丑
    ZhiYin             // 寅
    ZhiMao             // 卯
    ZhiChen            // 辰
    ZhiSi              // 巳
    ZhiWu              // 午
    ZhiWei             // 未
    ZhiShen            // 申
    ZhiYou             // 酉
    ZhiXu              // 戌
    ZhiHai             // 亥
)

// Branch 地支实体
type Branch struct {
    Zhi         Zhi
    MainElement Phase    // 主气：地支本气
    SubElements []Phase  // 余气：地支藏气
    Nature      Nature   // 阴阳属性
    Energy      uint8    // 能量级别（0-100）
}

// DiZhi 地支系统
type DiZhi struct {
    mu       sync.RWMutex
    branches map[Zhi]*Branch
    current  Zhi
    
    // 关联系统
    tianGan  *TianGan
    wuXing   *WuXing
    
    // 改进周期控制
    cycleManager struct {
        sync.RWMutex
        time      time.Duration
        last      time.Time
        schedule  map[time.Weekday][]CycleEvent
        active    bool
    }
    
    // 状态管理
    state      *state.StateManager
    ctx        *core.DaoContext
    metrics    *Metrics
    done       chan struct{}
}

// 新增周期事件
type CycleEvent struct {
    Time     time.Time
    Branch   Zhi
    Action   CycleAction
    Duration time.Duration
}

// NewDiZhi 创建地支系统
func NewDiZhi(ctx *core.DaoContext, tg *TianGan, wx *WuXing) *DiZhi {
    dz := &DiZhi{
        branches:  make(map[Zhi]*Branch),
        tianGan:   tg,
        wuXing:    wx,
        ctx:       ctx,
        cycleTime: time.Hour * 2,
        done:      make(chan struct{}),
    }
    
    dz.initBranches()
    return dz
}

// initBranches 初始化地支配置
func (dz *DiZhi) initBranches() {
    // 定义地支配置
    configs := map[Zhi]struct {
        main    Phase
        sub     []Phase
        nature  Nature
    }{
        ZhiZi:   {PhaseWater, []Phase{PhaseWater}, NatureYang},
        ZhiChou: {PhaseEarth, []Phase{PhaseEarth, PhaseMetal, PhaseWater}, NatureYin},
        ZhiYin:  {PhaseWood, []Phase{PhaseWood, PhaseFire, PhaseEarth}, NatureYang},
        ZhiMao:  {PhaseWood, []Phase{PhaseWood}, NatureYin},
        ZhiChen: {PhaseEarth, []Phase{PhaseEarth, PhaseWater, PhaseWood}, NatureYang},
        ZhiSi:   {PhaseFire, []Phase{PhaseFire, PhaseEarth, PhaseMetal}, NatureYin},
        ZhiWu:   {PhaseFire, []Phase{PhaseFire}, NatureYang},
        ZhiWei:  {PhaseEarth, []Phase{PhaseEarth, PhaseFire, PhaseWood}, NatureYin},
        ZhiShen: {PhaseMetal, []Phase{PhaseMetal, PhaseWater, PhaseEarth}, NatureYang},
        ZhiYou:  {PhaseMetal, []Phase{PhaseMetal}, NatureYin},
        ZhiXu:   {PhaseEarth, []Phase{PhaseEarth, PhaseFire, PhaseMetal}, NatureYang},
        ZhiHai:  {PhaseWater, []Phase{PhaseWater, PhaseWood}, NatureYin},
    }
    
    // 初始化每个地支
    for zhi, config := range configs {
        dz.branches[zhi] = &Branch{
            Zhi:         zhi,
            MainElement: config.main,
            SubElements: config.sub,
            Nature:      config.nature,
            Energy:      50, // 初始能量
        }
    }
    
    dz.current = ZhiZi // 默认从子开始
}

// Start 启动地支系统
func (dz *DiZhi) Start() error {
    dz.mu.Lock()
    if dz.running {
        dz.mu.Unlock()
        return errors.New("地支系统已在运行")
    }
    dz.running = true
    dz.lastCycle = time.Now()
    dz.mu.Unlock()
    
    go dz.runCycle()
    return nil
}

// runCycle 运行地支周期
func (dz *DiZhi) runCycle() {
    ticker := time.NewTicker(dz.cycleTime)
    defer ticker.Stop()
    
    for {
        select {
        case <-dz.done:
            return
        case <-ticker.C:
            dz.cycle()
        }
    }
}

// cycle 执行一次地支循环
func (dz *DiZhi) cycle() {
    dz.mu.Lock()
    defer dz.mu.Unlock()
    
    // 能量转换
    current := dz.branches[dz.current]
    next := dz.branches[(dz.current+1)%12]
    
    // 计算能量传递
    transfer := current.Energy / 10
    if transfer > 0 {
        current.Energy -= transfer
        next.Energy += transfer
        
        // 影响五行系统
        if dz.wuXing != nil {
            // 主气影响
            dz.wuXing.AdjustElement(next.MainElement, int8(transfer))
            
            // 余气影响
            for _, subElement := range next.SubElements {
                dz.wuXing.AdjustElement(subElement, int8(transfer/2))
            }
        }
    }
    
    // 更新当前地支
    dz.current = (dz.current + 1) % 12
    dz.lastCycle = time.Now()
}

// GetCurrent 获取当前地支信息
func (dz *DiZhi) GetCurrent() *Branch {
    dz.mu.RLock()
    defer dz.mu.RUnlock()
    
    return &Branch{
        Zhi:         dz.current,
        MainElement: dz.branches[dz.current].MainElement,
        SubElements: dz.branches[dz.current].SubElements,
        Nature:      dz.branches[dz.current].Nature,
        Energy:      dz.branches[dz.current].Energy,
    }
}

// AdjustEnergy 调整地支能量
func (dz *DiZhi) AdjustEnergy(zhi Zhi, delta int8) error {
    dz.mu.Lock()
    defer dz.mu.Unlock()
    
    branch, exists := dz.branches[zhi]
    if !exists {
        return ErrInvalidZhi
    }
    
    newEnergy := int16(branch.Energy) + int16(delta)
    if newEnergy < 0 {
        newEnergy = 0
    } else if newEnergy > 100 {
        newEnergy = 100
    }
    
    branch.Energy = uint8(newEnergy)
    return nil
}

// GetOpposite 获取地支对冲
func (dz *DiZhi) GetOpposite(zhi Zhi) Zhi {
    return (zhi + 6) % 12
}

// GetTriple 获取地支三合
func (dz *DiZhi) GetTriple(zhi Zhi) []Zhi {
    return []Zhi{
        zhi,
        (zhi + 4) % 12,
        (zhi + 8) % 12,
    }
}

// Stop 停止地支系统
func (dz *DiZhi) Stop() {
    dz.mu.Lock()
    if dz.running {
        dz.running = false
        close(dz.done)
    }
    dz.mu.Unlock()
}

// IsRunning 检查系统是否运行中
func (dz *DiZhi) IsRunning() bool {
    dz.mu.RLock()
    defer dz.mu.RUnlock()
    return dz.running
}
