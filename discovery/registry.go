// discovery/registry.go
package discovery

type ServiceRegistry struct {
    services  map[string]*Service
    watcher   *ServiceWatcher
    balancer  *LoadBalancer
    health    *HealthChecker
}

// discovery/resolver.go
type ServiceResolver struct {
    registry  *ServiceRegistry
    cache     *cache.Cache
    policy    ResolvePolicy
    metrics   *ResolverMetrics
}
