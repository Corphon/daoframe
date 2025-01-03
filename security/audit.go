//security/audit.go
package security

import (
    "context"
    "time"
)

// Auditor 审计服务
type Auditor struct {
    // 存储接口
    store      AuditStore
    // 事件过滤器
    filters    []AuditFilter
    // 事件处理器
    handlers   []AuditHandler
    // 监控指标
    metrics    *AuditMetrics
}

// AuditEvent 审计事件
type AuditEvent struct {
    ID          string
    Type        string
    Principal   *Principal
    Resource    string
    Action      string
    Result      string
    Error       error
    Metadata    map[string]interface{}
    Timestamp   time.Time
    Source      string
}

// AuditStore 审计存储接口
type AuditStore interface {
    SaveEvent(ctx context.Context, event *AuditEvent) error
    QueryEvents(ctx context.Context, filter *AuditFilter) ([]*AuditEvent, error)
    GetEventByID(ctx context.Context, id string) (*AuditEvent, error)
    DeleteEvents(ctx context.Context, filter *AuditFilter) error
}

// Log 记录审计日志
func (a *Auditor) Log(ctx context.Context, event *AuditEvent) error {
    // 应用过滤器
    for _, filter := range a.filters {
        if !filter.ShouldLog(event) {
            return nil
        }
    }

    // 丰富事件信息
    event.ID = generateEventID()
    event.Source = "security-service"
    if event.Timestamp.IsZero() {
        event.Timestamp = time.Now()
    }

    // 存储事件
    if err := a.store.SaveEvent(ctx, event); err != nil {
        a.metrics.LogErrors.Inc()
        return err
    }

    // 处理事件
    for _, handler := range a.handlers {
        if err := handler.HandleEvent(ctx, event); err != nil {
            a.metrics.HandlerErrors.Inc()
        }
    }

    a.metrics.EventsLogged.Inc()
    return nil
}
