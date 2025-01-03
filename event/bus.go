// event/bus.go
package event

type EventBus struct {
    handlers   map[string][]EventHandler
    middleware []EventMiddleware
    queue      *async.Queue[Event]
    metrics    *EventMetrics
}

// event/dispatcher.go
type EventDispatcher struct {
    bus       *EventBus
    filters   []EventFilter
    publisher *EventPublisher
    consumer  *EventConsumer
}

// event/handler.go
type EventHandler interface {
    Handle(ctx context.Context, event Event) error
    Priority() int
}
