//api/handler.go
package api

// Handler API处理器接口
type Handler interface {
    Handle(ctx *Context) error
    Validate(ctx *Context) error
}

// HandlerFunc 处理器函数类型
type HandlerFunc func(*Context) error

// HandlerRegistry 处理器注册表
type HandlerRegistry struct {
    handlers map[string]Handler
    mu       sync.RWMutex
}

// BaseHandler 基础处理器
type BaseHandler struct {
    validators []Validator
    logger     Logger
    metrics    *HandlerMetrics
}

func (h *BaseHandler) Validate(ctx *Context) error {
    for _, validator := range h.validators {
        if err := validator.Validate(ctx); err != nil {
            return &Error{
                Code:    http.StatusBadRequest,
                Message: "Validation failed",
                Details: map[string]string{
                    "field":   validator.Field(),
                    "reason":  err.Error(),
                },
            }
        }
    }
    return nil
}
