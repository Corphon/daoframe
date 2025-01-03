//security/access.go
package security

import (
    "context"
    "sync"
)

// AccessControl 访问控制服务
type AccessControl struct {
    // 策略存储
    policyStore PolicyStore
    // 策略评估器
    evaluator   *PolicyEvaluator
    // 策略缓存
    cache       *PolicyCache
    // 审计日志
    auditor     *Auditor
    // 监控指标
    metrics     *AccessMetrics
}

// PolicyStore 策略存储接口
type PolicyStore interface {
    GetPolicy(ctx context.Context, id string) (*Policy, error)
    ListPolicies(ctx context.Context, filter *PolicyFilter) ([]*Policy, error)
    CreatePolicy(ctx context.Context, policy *Policy) error
    UpdatePolicy(ctx context.Context, policy *Policy) error
    DeletePolicy(ctx context.Context, id string) error
}

// Policy 访问策略
type Policy struct {
    ID          string
    Name        string
    Description string
    Effect      EffectType
    Principals  []string
    Resources   []string
    Actions     []string
    Conditions  []Condition
    Priority    int
    Version     int64
    Created     time.Time
    Modified    time.Time
}

// Check 检查访问权限
func (ac *AccessControl) Check(ctx context.Context, principal *Principal, resource string, action string) (bool, error) {
    // 创建决策请求
    request := &AccessRequest{
        Principal: principal,
        Resource:  resource,
        Action:    action,
        Context:   ctx,
    }

    // 检查缓存
    if decision, found := ac.cache.Get(request); found {
        ac.metrics.CacheHits.Inc()
        return decision.Allowed, nil
    }
    ac.metrics.CacheMisses.Inc()

    // 评估策略
    decision, err := ac.evaluator.Evaluate(ctx, request)
    if err != nil {
        ac.metrics.EvaluationErrors.Inc()
        return false, err
    }

    // 记录审计日志
    ac.auditor.Log(ctx, &AuditEvent{
        Principal: principal,
        Resource:  resource,
        Action:    action,
        Decision:  decision,
        Timestamp: time.Now(),
    })

    // 更新缓存
    ac.cache.Set(request, decision)

    return decision.Allowed, nil
}
