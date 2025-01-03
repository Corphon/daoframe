// api/router.go
package api

type Router struct {
    routes     map[string][]Route
    middleware []Middleware
    handlers   *HandlerRegistry
    metrics    *RouterMetrics
}

// api/handler.go
type Handler interface {
    Handle(ctx *Context) error
    Validate(ctx *Context) error
}

// api/middleware.go
type Middleware interface {
    Before(ctx *Context) error
    After(ctx *Context) error
}
