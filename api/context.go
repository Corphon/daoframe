//api/context.go
package api

import (
    "context"
    "encoding/json"
    "net/http"
    "sync"
    
    "github.com/Corphon/daoframe/tools/cache"
)

// Context API上下文
type Context struct {
    context.Context
    Request        *http.Request
    ResponseWriter http.ResponseWriter
    Params         map[string]string
    route         *Route
    handler       Handler
    middleware    []Middleware
    
    // 请求相关
    requestID     string
    startTime     time.Time
    timeout       time.Duration
    
    // 数据存储
    store         sync.Map
    cache         *cache.Cache
    
    // 状态标记
    wrote         bool
    aborted      bool
    
    // 性能追踪
    tracer       *Tracer
    spans        []*Span
}

// NewContext 创建新的上下文
func NewContext(w http.ResponseWriter, r *http.Request, opts ...ContextOption) *Context {
    ctx := &Context{
        Context:        r.Context(),
        Request:        r,
        ResponseWriter: w,
        Params:        make(map[string]string),
        startTime:     time.Now(),
        cache:         cache.New(cache.DefaultExpiration),
    }
    
    // 应用选项
    for _, opt := range opts {
        opt(ctx)
    }
    
    // 初始化追踪
    ctx.tracer = NewTracer(ctx.requestID)
    ctx.startSpan("request")
    
    return ctx
}

// JSON 返回JSON响应
func (c *Context) JSON(code int, data interface{}) error {
    if c.wrote {
        return errors.New("response already written")
    }
    
    resp := Response{
        Code:      code,
        Data:      data,
        RequestID: c.requestID,
        Timestamp: time.Now(),
    }
    
    c.ResponseWriter.Header().Set("Content-Type", ContentTypeJSON)
    c.ResponseWriter.WriteHeader(code)
    
    if err := json.NewEncoder(c.ResponseWriter).Encode(resp); err != nil {
        return err
    }
    
    c.wrote = true
    return nil
}

// Error 返回错误响应
func (c *Context) Error(err error) {
    var apiErr *Error
    if errors.As(err, &apiErr) {
        c.JSON(apiErr.Code, apiErr)
        return
    }
    
    // 包装普通错误
    apiErr = &Error{
        Code:    http.StatusInternalServerError,
        Message: err.Error(),
    }
    c.JSON(apiErr.Code, apiErr)
}
