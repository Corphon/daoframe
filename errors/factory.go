//errors/factory.go
package errors

import (
    "runtime"
    "time"
)

// ErrorFactory 错误工厂
type ErrorFactory struct {
    defaultCategory ErrorCategory
    defaultLevel   ErrorLevel
    withStack      bool
}

// NewErrorFactory 创建错误工厂
func NewErrorFactory(category ErrorCategory) *ErrorFactory {
    return &ErrorFactory{
        defaultCategory: category,
        defaultLevel:   LevelError,
        withStack:      true,
    }
}

// New 创建新错误
func (f *ErrorFactory) New(code ErrorCode, msg string) *BaseError {
    err := &BaseError{
        Code:     code,
        Category: f.defaultCategory,
        Level:    f.defaultLevel,
        Message:  msg,
        Context: ErrorContext{
            Time:    time.Now(),
            Details: make(map[string]interface{}),
        },
    }
    
    if f.withStack {
        err.Context.Stack = f.getStack()
    }
    
    return err
}

// WithCause 添加原因
func (f *ErrorFactory) WithCause(code ErrorCode, msg string, cause error) *BaseError {
    err := f.New(code, msg)
    err.Cause = cause
    return err
}
