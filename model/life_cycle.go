// model/life_cycle.go

package model

import (
    "sync"
    "time"
    "errors"
    
    "github.com/Corphon/daoframe/core"
)

var (
    ErrCycleInvalid    = errors.New("无效的生命周期")
    ErrStateTransition = errors.New("状态转换错误")
    ErrCycleLocked     = errors.New("生命周期已锁定")
)

// LifeStage 生命阶段
type LifeStage uint8

const (
    StageVoid   LifeStage = iota // 虚无
    StageOrigin                   // 混沌
    StageBirth                    // 生成
    StageGrowth                   // 生长
    StagePeak                     // 极盛
    StageDecline                  // 衰退
    StageEnd                      // 终结
    StageReturn                   // 归一
)

// Element 元素特性
type LifeElement struct {
    Phase     Phase     // 五行属性
    Nature    Nature    // 阴阳属性
    Energy    uint8     // 能量值
    Vitality  uint8     // 生命力
}

// LifeEntity 生命实体
type LifeEntity struct {
    ID        string
    Stage     LifeStage
    Elements  []LifeElement
    Birth     time.Time
    LastCycle time.Time
    Duration  time.Duration
}

// LifeCycle 生命周期系统
type LifeCycle struct {
    mu       sync.RWMutex
    entities map[string]*LifeEntity
    observers   []LifeCycleObserver    // 新增：观察者列表
    entityLocks map[string]*sync.RWMutex  // 新增：实体级别锁
    lockShards  []*sync.RWMutex          // 新增：分片锁
    
    // 关联系统
    wuXing   *WuXing
    tianGan  *TianGan
    diZhi    *DiZhi
    
    // 周期控制
    ctx      *core.DaoContext
    running  bool
    done     chan struct{}
}

// NewLifeCycle 创建生命周期系统
func NewLifeCycle(ctx *core.DaoContext, wx *WuXing, tg *TianGan, dz *DiZhi) *LifeCycle {
    return &LifeCycle{
        entities: make(map[string]*LifeEntity),
        observers:   make([]LifeCycleObserver, 0),
        entityLocks: make(map[string]*sync.RWMutex),
        lockShards:  make([]*sync.RWMutex, 32), // 32个分片锁
        wuXing:   wx,
        tianGan:  tg,
        diZhi:    dz,
        ctx:      ctx,
        done:     make(chan struct{}),
    }
     // 初始化分片锁
    for i := range lc.lockShards {
        lc.lockShards[i] = &sync.RWMutex{}
    }
    
    return lc
}

// CreateEntity 创建生命实体
func (lc *LifeCycle) CreateEntity(id string) (*LifeEntity, error) {
    lc.mu.Lock()
    defer lc.mu.Unlock()

    if _, exists := lc.entities[id]; exists {
        return nil, errors.New("实体已存在")
    }

    // 创建新实体
    entity := &LifeEntity{
        ID:        id,
        Stage:     StageVoid,
        Elements:  make([]LifeElement, 0),
        Birth:     time.Now(),
        LastCycle: time.Now(),
        Duration:  0,
    }

    // 基于当前天干地支初始化元素
    if lc.tianGan != nil && lc.diZhi != nil {
        currentGan := lc.tianGan.GetCurrentGan()
        currentZhi := lc.diZhi.GetCurrent()
        
        // 添加主元素
        entity.Elements = append(entity.Elements, LifeElement{
            Phase:    lc.tianGan.GetGanElement(currentGan),
            Nature:   lc.tianGan.GetGanNature(currentGan),
            Energy:   50,
            Vitality: 100,
        })
        
        // 添加地支元素
        entity.Elements = append(entity.Elements, LifeElement{
            Phase:    currentZhi.MainElement,
            Nature:   currentZhi.Nature,
            Energy:   50,
            Vitality: 100,
        })
    }

    lc.entities[id] = entity
    return entity, nil
}

// Start 启动生命周期系统
func (lc *LifeCycle) Start() error {
    lc.mu.Lock()
    if lc.running {
        lc.mu.Unlock()
        return errors.New("系统已在运行")
    }
    lc.running = true
    lc.mu.Unlock()

    go lc.runCycles()
    return nil
}

// runCycles 运行生命周期
func (lc *LifeCycle) runCycles() {
    ticker := time.NewTicker(time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-lc.done:
            return
        case <-ticker.C:
            lc.processCycle()
        }
    }
}

// processCycle 处理生命周期
func (lc *LifeCycle) processCycle() {
    lc.mu.Lock()
    defer lc.mu.Unlock()

    now := time.Now()
    for _, entity := range lc.entities {
        // 更新实体状态
        lc.updateEntityState(entity, now)
        
        // 处理元素变化
        lc.processElementChanges(entity)
        
        entity.LastCycle = now
    }
}

// updateEntityState 更新实体状态
func (lc *LifeCycle) updateEntityState(entity *LifeEntity, now time.Time) {
    age := now.Sub(entity.Birth)
    entity.Duration = age
    oldStage := entity.Stage

    // 基于年龄和生命力确定阶段
    totalVitality := lc.calculateTotalVitality(entity)
    
    switch {
    case totalVitality == 0:
        entity.Stage = StageReturn
    case totalVitality < 20:
        entity.Stage = StageEnd
    case totalVitality < 40:
        entity.Stage = StageDecline
    case totalVitality < 60:
        entity.Stage = StageGrowth
    case totalVitality < 80:
        entity.Stage = StagePeak
    default:
        entity.Stage = StageBirth
    }
    // 如果状态发生变化，通知观察者
    if oldStage != entity.Stage {
        event := LifeEvent{
            EntityID:   entity.ID,
            OldStage:   oldStage,
            NewStage:   entity.Stage,
            TimeStamp:  now,
        }
        
        for _, observer := range lc.observers {
            go observer.OnStateChange(event)
        }
    }
}

// calculateTotalVitality 计算总生命力
func (lc *LifeCycle) calculateTotalVitality(entity *LifeEntity) uint8 {
    if len(entity.Elements) == 0 {
        return 0
    }

    var total uint16
    for _, elem := range entity.Elements {
        total += uint16(elem.Vitality)
    }
    
    return uint8(total / uint16(len(entity.Elements)))
}

// processElementChanges 处理元素变化
func (lc *LifeCycle) processElementChanges(entity *LifeEntity) {
    for i := range entity.Elements {
        elem := &entity.Elements[i]
        
        // 基于五行相生相克调整能量
        if lc.wuXing != nil {
            strength, _ := lc.wuXing.GetElementStrength(elem.Phase)
            energyDelta := (int8(strength) - int8(elem.Energy)) / 10
            
            newEnergy := int16(elem.Energy) + int16(energyDelta)
            if newEnergy < 0 {
                newEnergy = 0
            } else if newEnergy > 100 {
                newEnergy = 100
            }
            elem.Energy = uint8(newEnergy)
        }

        // 生命力随时间自然衰减
        if elem.Vitality > 0 {
            elem.Vitality--
        }
    }
}

// GetEntity 获取实体信息
func (lc *LifeCycle) GetEntity(id string) (*LifeEntity, error) {
    lc.mu.RLock()
    defer lc.mu.RUnlock()

    entity, exists := lc.entities[id]
    if !exists {
        return nil, errors.New("实体不存在")
    }

    // 返回副本
    return &LifeEntity{
        ID:        entity.ID,
        Stage:     entity.Stage,
        Elements:  append([]LifeElement{}, entity.Elements...),
        Birth:     entity.Birth,
        LastCycle: entity.LastCycle,
        Duration:  entity.Duration,
    }, nil
}

// Stop 停止生命周期系统
func (lc *LifeCycle) Stop() {
    lc.mu.Lock()
    if lc.running {
        lc.running = false
        close(lc.done)
    }
    lc.mu.Unlock()
}
