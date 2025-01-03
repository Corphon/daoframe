//api/middleware.go
package api

// Middleware 中间件接口
type Middleware interface {
    Process(ctx *Context, next HandlerFunc) error
}

// LoggerMiddleware 日志中间件
type LoggerMiddleware struct {
    logger Logger
}

func (m *LoggerMiddleware) Process(ctx *Context, next HandlerFunc) error {
    // 开始计时
    start := time.Now()
    
    // 记录请求信息
    m.logger.Info("API Request",
        "method", ctx.Request.Method,
        "path", ctx.Request.URL.Path,
        "request_id", ctx.requestID,
    )
    
    // 执行下一个处理器
    err := next(ctx)
    
    // 记录响应信息
    duration := time.Since(start)
    m.logger.Info("API Response",
        "method", ctx.Request.Method,
        "path", ctx.Request.URL.Path,
        "duration", duration,
        "status", ctx.ResponseWriter.Status(),
        "error", err,
    )
    
    return err
}

// RecoveryMiddleware 恢复中间件
type RecoveryMiddleware struct {
    logger Logger
}

func (m *RecoveryMiddleware) Process(ctx *Context, next HandlerFunc) error {
    defer func() {
        if r := recover(); r != nil {
            err := fmt.Errorf("panic recovered: %v", r)
            stack := debug.Stack()
            
            m.logger.Error("API Panic",
                "error", err,
                "stack", string(stack),
            )
            
            ctx.Error(&Error{
                Code:    http.StatusInternalServerError,
                Message: "Internal Server Error",
                Stack:   string(stack),
            })
        }
    }()
    
    return next(ctx)
}
