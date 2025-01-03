//discovery/health.go
package discovery

// HealthChecker 健康检查器
type HealthChecker struct {
    registry   *ServiceRegistry
    checkers   map[string]HealthCheck
    intervals  map[string]time.Duration
    metrics    *HealthMetrics
    
    ctx        context.Context
    cancel     context.CancelFunc
    wg         sync.WaitGroup
}

// HealthCheck 健康检查接口
type HealthCheck interface {
    Check(ctx context.Context, instance *ServiceInstance) *HealthResult
    Name() string
}

// HealthResult 健康检查结果
type HealthResult struct {
    Status    ServiceStatus
    Error     error
    Details   map[string]interface{}
    Timestamp time.Time
}

// HTTPHealthCheck HTTP健康检查
type HTTPHealthCheck struct {
    client    *http.Client
    path      string
    timeout   time.Duration
}

func (h *HTTPHealthCheck) Check(ctx context.Context, instance *ServiceInstance) *HealthResult {
    url := fmt.Sprintf("%s%s", instance.Endpoint, h.path)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return &HealthResult{
            Status:    StatusDown,
            Error:     err,
            Timestamp: time.Now(),
        }
    }

    resp, err := h.client.Do(req)
    if err != nil {
        return &HealthResult{
            Status:    StatusDown,
            Error:     err,
            Timestamp: time.Now(),
        }
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        return &HealthResult{
            Status:    StatusUp,
            Timestamp: time.Now(),
        }
    }

    return &HealthResult{
        Status:    StatusDown,
        Error:     fmt.Errorf("unhealthy status code: %d", resp.StatusCode),
        Timestamp: time.Now(),
    }
}
