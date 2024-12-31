// core/universe.go

package core

import (
    "context"
    "sync"
    "time"
)

// Phase 表示宇宙演化阶段
type Phase uint8

const (
    PhaseWuJi   Phase = iota // 无极：混沌未分
    PhaseTaiJi              // 太极：道生一
    PhaseYinYang            // 阴阳：一生二
    PhaseTriad              // 三态：二生三
    PhaseWuXing             // 五行：三生五
    PhaseWanWu              // 万物：五生万
)

// Element 表示基本元素
type Element struct {
    mu       sync.RWMutex
    yin      float64  // 阴性比例 0-1
    yang     float64  // 阳性比例 0-1
    energy   float64  // 能量水平
    phase    Phase    // 所处阶段
}

// Universe 表示宇宙整体
type Universe struct {
    mu          sync.RWMutex
    ctx         context.Context
    phase       Phase
    elements    []*Element
    triad       *Triad      // 三态系统
    wuxing      *WuXing     // 五行系统
    observers   []Observer
    done        chan struct{}
}

// NewUniverse 创建新的宇宙实例，从无极开始
func NewUniverse(ctx context.Context) *Universe {
    return &Universe{
        ctx:       ctx,
        phase:     PhaseWuJi,
        elements:  make([]*Element, 0),
        observers: make([]Observer, 0),
        done:      make(chan struct{}),
    }
}

// Evolve 开始宇宙演化
func (u *Universe) Evolve() error {
    // 1. 无极 -> 太极（道生一）
    if err := u.generateTaiJi(); err != nil {
        return err
    }

    // 2. 太极 -> 阴阳（一生二）
    if err := u.generateYinYang(); err != nil {
        return err
    }

    // 3. 阴阳 -> 三态（二生三）
    if err := u.generateTriad(); err != nil {
        return err
    }

    // 4. 三态 -> 五行（三生五）
    if err := u.generateWuXing(); err != nil {
        return err
    }

    // 5. 五行生万物
    go u.generateWanWu()

    return nil
}

// generateTaiJi 实现"道生一"：从无极生成太极
func (u *Universe) generateTaiJi() error {
    u.mu.Lock()
    defer u.mu.Unlock()

    if u.phase != PhaseWuJi {
        return ErrInvalidPhase
    }

    // 创建原始元素，能量完全平衡
    primordial := &Element{
        yin:    0.5,
        yang:   0.5,
        energy: 1.0,
        phase:  PhaseTaiJi,
    }
    u.elements = append(u.elements, primordial)
    u.phase = PhaseTaiJi
    
    return nil
}

// generateYinYang 实现"一生二"：太极分化为阴阳
func (u *Universe) generateYinYang() error {
    u.mu.Lock()
    defer u.mu.Unlock()

    if u.phase != PhaseTaiJi {
        return ErrInvalidPhase
    }

    // 从原始元素分化出阴阳
    original := u.elements[0]
    
    // 创建阴元素
    yin := &Element{
        yin:    0.8,
        yang:   0.2,
        energy: original.energy / 2,
        phase:  PhaseYinYang,
    }
    
    // 创建阳元素
    yang := &Element{
        yin:    0.2,
        yang:   0.8,
        energy: original.energy / 2,
        phase:  PhaseYinYang,
    }

    u.elements = append(u.elements, yin, yang)
    u.phase = PhaseYinYang
    
    return nil
}

// Triad 三态系统
type Triad struct {
    minorYin  *Element // 少阴
    minorYang *Element // 少阳
    balance   *Element // 中和
}

// generateTriad 实现"二生三"：阴阳相互作用产生三态
func (u *Universe) generateTriad() error {
    u.mu.Lock()
    defer u.mu.Unlock()

    if u.phase != PhaseYinYang {
        return ErrInvalidPhase
    }

    // 从阴阳相互作用生成三态
    u.triad = &Triad{
        minorYin: &Element{
            yin:    0.7,
            yang:   0.3,
            energy: 1.0,
            phase:  PhaseTriad,
        },
        minorYang: &Element{
            yin:    0.3,
            yang:   0.7,
            energy: 1.0,
            phase:  PhaseTriad,
        },
        balance: &Element{
            yin:    0.5,
            yang:   0.5,
            energy: 1.0,
            phase:  PhaseTriad,
        },
    }

    u.elements = append(u.elements, 
        u.triad.minorYin, 
        u.triad.minorYang, 
        u.triad.balance,
    )
    u.phase = PhaseTriad
    
    return nil
}

// generateWuXing 实现"三生五"：三态演化为五行
func (u *Universe) generateWuXing() error {
    u.mu.Lock()
    defer u.mu.Unlock()

    if u.phase != PhaseTriad {
        return ErrInvalidPhase
    }

    // 从三态演化出五行
    wood := &Element{yin: 0.4, yang: 0.6, energy: 1.0, phase: PhaseWuXing}
    fire := &Element{yin: 0.2, yang: 0.8, energy: 1.0, phase: PhaseWuXing}
    earth := &Element{yin: 0.5, yang: 0.5, energy: 1.0, phase: PhaseWuXing}
    metal := &Element{yin: 0.6, yang: 0.4, energy: 1.0, phase: PhaseWuXing}
    water := &Element{yin: 0.8, yang: 0.2, energy: 1.0, phase: PhaseWuXing}

    u.wuxing = NewWuXing([]*Element{wood, fire, earth, metal, water})
    u.elements = append(u.elements, wood, fire, earth, metal, water)
    u.phase = PhaseWuXing

    return nil
}

// generateWanWu 实现"五生万物"：持续的演化过程
func (u *Universe) generateWanWu() {
    ticker := time.NewTicker(time.Millisecond * 100)
    defer ticker.Stop()

    for {
        select {
        case <-u.done:
            return
        case <-ticker.C:
            u.mu.Lock()
            if u.phase == PhaseWuXing {
                // 通过五行相互作用生成新元素
                if newElement := u.wuxing.Interact(); newElement != nil {
                    newElement.phase = PhaseWanWu
                    u.elements = append(u.elements, newElement)
                    // 通知观察者
                    u.notifyObservers(EventNewElement, newElement)
                }
            }
            u.mu.Unlock()
        }
    }
}

// 观察者相关定义
type EventType uint8

const (
    EventNewElement EventType = iota
    EventPhaseChange
    EventInteraction
)

type Observer interface {
    OnEvent(eventType EventType, data interface{})
}

func (u *Universe) AddObserver(observer Observer) {
    u.mu.Lock()
    defer u.mu.Unlock()
    u.observers = append(u.observers, observer)
}

func (u *Universe) notifyObservers(eventType EventType, data interface{}) {
    for _, observer := range u.observers {
        observer.OnEvent(eventType, data)
    }
}
