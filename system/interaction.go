// system/interaction.go

package system

import (
    "sync"
    "time"
    
    "github.com/Corphon/daoframe/core"
    "github.com/Corphon/daoframe/model/basic"
)

// InteractionType 交互类型
type InteractionType uint8

const (
    // 基础交互类型
    BaguaInteraction InteractionType = iota  // 八卦间交互
    ElementInteraction                       // 五行间交互
    TemporalInteraction                      // 时序交互
    
    // 复合交互类型
    BaguaElement                             // 八卦与五行
    BaguaTemporal                            // 八卦与时序
    ElementTemporal                          // 五行与时序
)

// Interaction 交互事件
type Interaction struct {
    Type        InteractionType
    Source      interface{}           // 源对象
    Target      interface{}           // 目标对象
    Strength    float64              // 交互强度
    Duration    time.Duration        // 持续时间
    Timestamp   time.Time            // 发生时间
}

// InteractionEffect 交互效果
type InteractionEffect struct {
    EnergyDelta     float64              // 能量变化
    AttributeChanges map[string]float64   // 属性变化
    StateChanges    map[string]interface{} // 状态变化
}

// InteractionSystem 交互系统
type InteractionSystem struct {
    mu              sync.RWMutex
    ctx             *core.DaoContext
    
    // 基础系统引用
    bagua           *basic.BaGua
    wuXing          *basic.WuXing
    timeSystem      *basic.TimeSystem
    
    // 交互规则
    rules           map[InteractionType][]InteractionRule
    
    // 交互历史
    history         []Interaction
    
    // 效果缓存
    effectCache     map[string]*InteractionEffect
    
    // 观察者
    observers       []InteractionObserver
    
    // 控制通道
    done            chan struct{}
}

// InteractionRule 交互规则
type InteractionRule struct {
    Condition    func(*Interaction) bool           // 触发条件
    Effect       func(*Interaction) *InteractionEffect // 效果计算
    Priority     int                               // 规则优先级
}

// NewInteractionSystem 创建交互系统
func NewInteractionSystem(ctx *core.DaoContext, bagua *basic.BaGua, 
    wuXing *basic.WuXing, timeSystem *basic.TimeSystem) *InteractionSystem {
    
    is := &InteractionSystem{
        ctx:         ctx,
        bagua:       bagua,
        wuXing:      wuXing,
        timeSystem:  timeSystem,
        rules:       make(map[InteractionType][]InteractionRule),
        history:     make([]Interaction, 0),
        effectCache: make(map[string]*InteractionEffect),
        observers:   make([]InteractionObserver, 0),
        done:        make(chan struct{}),
    }
    
    // 初始化基础规则
    is.initializeRules()
    
    return is
}

// RegisterRule 注册交互规则
func (is *InteractionSystem) RegisterRule(typ InteractionType, rule InteractionRule) {
    is.mu.Lock()
    defer is.mu.Unlock()
    
    is.rules[typ] = append(is.rules[typ], rule)
    
    // 按优先级排序
    sort.Slice(is.rules[typ], func(i, j int) bool {
        return is.rules[typ][i].Priority > is.rules[typ][j].Priority
    })
}

// ProcessInteraction 处理交互
func (is *InteractionSystem) ProcessInteraction(interaction *Interaction) *InteractionEffect {
    is.mu.Lock()
    defer is.mu.Unlock()
    
    // 检查缓存
    cacheKey := is.generateCacheKey(interaction)
    if effect, exists := is.effectCache[cacheKey]; exists {
        return effect
    }
    
    // 应用规则
    effect := is.applyRules(interaction)
    
    // 记录历史
    is.history = append(is.history, *interaction)
    
    // 缓存结果
    is.effectCache[cacheKey] = effect
    
    // 通知观察者
    is.notifyObservers(interaction, effect)
    
    return effect
}

// initializeRules 初始化基础规则
func (is *InteractionSystem) initializeRules() {
    // 八卦交互规则
    is.RegisterRule(BaguaInteraction, InteractionRule{
        Condition: func(i *Interaction) bool {
            // 检查八卦交互条件
            return true
        },
        Effect: func(i *Interaction) *InteractionEffect {
            return is.calculateBaguaInteraction(i)
        },
        Priority: 10,
    })
    
    // 五行交互规则
    is.RegisterRule(ElementInteraction, InteractionRule{
        Condition: func(i *Interaction) bool {
            // 检查五行交互条件
            return true
        },
        Effect: func(i *Interaction) *InteractionEffect {
            return is.calculateElementInteraction(i)
        },
        Priority: 8,
    })
    
    // 时序交互规则
    is.RegisterRule(TemporalInteraction, InteractionRule{
        Condition: func(i *Interaction) bool {
            // 检查时序交互条件
            return true
        },
        Effect: func(i *Interaction) *InteractionEffect {
            return is.calculateTemporalInteraction(i)
        },
        Priority: 5,
    })
}

// calculateBaguaInteraction 计算八卦交互效果
func (is *InteractionSystem) calculateBaguaInteraction(i *Interaction) *InteractionEffect {
    effect := &InteractionEffect{
        EnergyDelta:     0,
        AttributeChanges: make(map[string]float64),
        StateChanges:    make(map[string]interface{}),
    }
    
    sourceTrigram := i.Source.(basic.Trigram)
    targetTrigram := i.Target.(basic.Trigram)
    
    // 计算能量变化
    effect.EnergyDelta = is.bagua.CalculateEnergyExchange(sourceTrigram, targetTrigram)
    
    // 计算属性变化
    attributes := is.bagua.GetTrigramAttributes(sourceTrigram, targetTrigram)
    for k, v := range attributes {
        effect.AttributeChanges[k] = v
    }
    
    return effect
}

// calculateElementInteraction 计算五行交互效果
func (is *InteractionSystem) calculateElementInteraction(i *Interaction) *InteractionEffect {
    effect := &InteractionEffect{
        EnergyDelta:     0,
        AttributeChanges: make(map[string]float64),
        StateChanges:    make(map[string]interface{}),
    }
    
    sourcePhase := i.Source.(basic.Phase)
    targetPhase := i.Target.(basic.Phase)
    
    // 计算五行相互作用
    relationship := is.wuXing.GetRelationship(sourcePhase, targetPhase)
    effect.EnergyDelta = is.wuXing.CalculateInteraction(sourcePhase, targetPhase, relationship)
    
    return effect
}

// calculateTemporalInteraction 计算时序交互效果
func (is *InteractionSystem) calculateTemporalInteraction(i *Interaction) *InteractionEffect {
    effect := &InteractionEffect{
        EnergyDelta:     0,
        AttributeChanges: make(map[string]float64),
        StateChanges:    make(map[string]interface{}),
    }
    
    // 计算时序影响
    temporal := is.timeSystem.GetCurrentCycle()
    effect.EnergyDelta = is.timeSystem.CalculateTemporalEffect(temporal)
    
    return effect
}

// AddObserver 添加观察者
func (is *InteractionSystem) AddObserver(observer InteractionObserver) {
    is.mu.Lock()
    defer is.mu.Unlock()
    is.observers = append(is.observers, observer)
}

// notifyObservers 通知观察者
func (is *InteractionSystem) notifyObservers(i *Interaction, effect *InteractionEffect) {
    for _, observer := range is.observers {
        observer.OnInteraction(i, effect)
    }
}

// generateCacheKey 生成缓存键
func (is *InteractionSystem) generateCacheKey(i *Interaction) string {
    return fmt.Sprintf("%d-%v-%v-%f-%s", 
        i.Type, i.Source, i.Target, i.Strength, i.Timestamp)
}

// Close 关闭交互系统
func (is *InteractionSystem) Close() {
    close(is.done)
}
