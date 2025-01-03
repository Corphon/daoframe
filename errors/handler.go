//errors/handler.go
package errors

// ErrorHandler 错误处理器接口
type ErrorHandler interface {
    Handle(error) error
    Recover() error
}

// DefaultErrorHandler 默认错误处理器
type DefaultErrorHandler struct {
    factory  *ErrorFactory
    handlers map[ErrorCategory][]ErrorProcessor
    logger   Logger
}

// ErrorProcessor 错误处理函数
type ErrorProcessor func(*BaseError) error

// NewDefaultErrorHandler 创建默认错误处理器
func NewDefaultErrorHandler(logger Logger) *DefaultErrorHandler {
    return &DefaultErrorHandler{
        factory:  NewErrorFactory(CategorySystem),
        handlers: make(map[ErrorCategory][]ErrorProcessor),
        logger:   logger,
    }
}

// RegisterProcessor 注册错误处理器
func (h *DefaultErrorHandler) RegisterProcessor(category ErrorCategory, processor ErrorProcessor) {
    if h.handlers[category] == nil {
        h.handlers[category] = make([]ErrorProcessor, 0)
    }
    h.handlers[category] = append(h.handlers[category], processor)
}
