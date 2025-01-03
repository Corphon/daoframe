//errors/types.go
package errors

import (
    "fmt"
    "encoding/json"
)

// ErrorCode 错误码类型
type ErrorCode int

// ErrorLevel 错误级别
type ErrorLevel uint8

const (
    LevelDebug ErrorLevel = iota
    LevelInfo
    LevelWarn
    LevelError
    LevelFatal
)

// ErrorCategory 错误类别
type ErrorCategory string

const (
    CategorySystem   ErrorCategory = "system"
    CategoryCore     ErrorCategory = "core"
    CategoryModel    ErrorCategory = "model"
    CategoryConfig   ErrorCategory = "config"
    CategoryRuntime  ErrorCategory = "runtime"
)

// BaseError 基础错误结构
type BaseError struct {
    Code     ErrorCode     `json:"code"`
    Category ErrorCategory `json:"category"`
    Level    ErrorLevel    `json:"level"`
    Message  string       `json:"message"`
    Cause    error        `json:"cause,omitempty"`
    Context  ErrorContext `json:"context,omitempty"`
}

// ErrorContext 错误上下文
type ErrorContext struct {
    Time    time.Time              `json:"time"`
    Stack   string                 `json:"stack,omitempty"`
    Details map[string]interface{} `json:"details,omitempty"`
}
