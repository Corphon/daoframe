// api/router.go
package api

import (
    "context"
    "net/http"
    "sync"
    
    "github.com/Corphon/daoframe/tools/async"
)

// Router API路由器
type Router struct {
    routes     map[string]*RouteGroup
    middleware []Middleware
    handlers   *HandlerRegistry
    
    // 配置
    config     *RouterConfig
    // 工作池
    pool       *async.WorkerPool
    // 限流器
    limiter    *RateLimiter
    // 监控
    metrics    *RouterMetrics
    
    mu         sync.RWMutex
}

// RouterConfig 路由器配置
type RouterConfig struct {
    // 基础路径
    BasePath string
    // 最大请求体大小
    MaxBodySize int64
    // 超时设置
    Timeout time.Duration
    // 工作池大小
    WorkerPoolSize int
    // 限流设置
    RateLimit RateLimitConfig
    // 跨域设置
    CORS CORSConfig
}

// Route 路由
type Route struct {
    Path        string
    Method      Method
    Handler     Handler
    Middleware  []Middleware
    Validators  []Validator
    Description string
    Tags        []string
}

// RouteGroup 路由组
type RouteGroup struct {
    prefix     string
    routes     []*Route
    middleware []Middleware
}

// Handle 处理请求
func (r *Router) Handle(w http.ResponseWriter, req *http.Request) {
    // 创建上下文
    ctx := NewContext(w, req)
    
    // 应用全局中间件
    handler := r.applyMiddleware(ctx)
    
    // 查找路由
    route, params, err := r.findRoute(req.Method, req.URL.Path)
    if err != nil {
        ctx.Error(err)
        return
    }
    
    // 设置路由参数
    ctx.Params = params
    ctx.route = route
    
    // 限流检查
    if err := r.limiter.Allow(ctx); err != nil {
        ctx.Error(err)
        return
    }
    
    // 提交到工作池处理
    r.pool.Submit(func() {
        defer r.recoverPanic(ctx)
        
        // 执行请求
        if err := handler(ctx); err != nil {
            ctx.Error(err)
        }
    })
}

// Group 创建路由组
func (r *Router) Group(prefix string, middleware ...Middleware) *RouteGroup {
    return &RouteGroup{
        prefix:     prefix,
        middleware: middleware,
    }
}
