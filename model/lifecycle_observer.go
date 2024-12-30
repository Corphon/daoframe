// model/lifecycle_observer.go

package model

import "time"

// 生命周期事件定义
type LifeEvent struct {
    EntityID    string
    OldStage    LifeStage
    NewStage    LifeStage
    TimeStamp   time.Time
}

// 观察者接口
type LifeCycleObserver interface {
    OnStateChange(event LifeEvent)
}

// 默认观察者实现
type DefaultLifeCycleObserver struct {
    // 可以添加日志记录器等
}

func (o *DefaultLifeCycleObserver) OnStateChange(event LifeEvent) {
    // 默认实现
}
