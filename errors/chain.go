//errors/chain.go
package errors

// ErrorChain 错误链
type ErrorChain struct {
    errors  []error
    context ErrorContext
}

// NewErrorChain 创建错误链
func NewErrorChain() *ErrorChain {
    return &ErrorChain{
        errors:  make([]error, 0),
        context: ErrorContext{
            Time:    time.Now(),
            Details: make(map[string]interface{}),
        },
    }
}

// Add 添加错误
func (c *ErrorChain) Add(err error) {
    if err != nil {
        c.errors = append(c.errors, err)
    }
}

// HasErrors 检查是否有错误
func (c *ErrorChain) HasErrors() bool {
    return len(c.errors) > 0
}
