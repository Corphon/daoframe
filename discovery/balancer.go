//discovery/balancer.go
package discovery

// LoadBalancer 负载均衡器
type LoadBalancer struct {
    strategies map[string]BalanceStrategy
    metrics    *BalancerMetrics
}

// BalanceStrategy 负载均衡策略接口
type BalanceStrategy interface {
    Select(instances []*ServiceInstance) (*ServiceInstance, error)
    Name() string
}

// RoundRobinStrategy 轮询策略
type RoundRobinStrategy struct {
    current uint64
}

func (s *RoundRobinStrategy) Select(instances []*ServiceInstance) (*ServiceInstance, error) {
    if len(instances) == 0 {
        return nil, errors.NoAvailableInstances
    }
    
    index := atomic.AddUint64(&s.current, 1) % uint64(len(instances))
    return instances[index], nil
}

// WeightedRandomStrategy 加权随机策略
type WeightedRandomStrategy struct {
    rand *rand.Rand
}

func (s *WeightedRandomStrategy) Select(instances []*ServiceInstance) (*ServiceInstance, error) {
    if len(instances) == 0 {
        return nil, errors.NoAvailableInstances
    }

    // 计算总权重
    totalWeight := 0
    for _, inst := range instances {
        totalWeight += inst.Weight
    }

    // 随机选择
    target := s.rand.Intn(totalWeight)
    current := 0
    
    for _, inst := range instances {
        current += inst.Weight
        if current > target {
            return inst, nil
        }
    }

    return instances[0], nil
}
