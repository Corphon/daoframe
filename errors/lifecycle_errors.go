// errors/lifecycle_errors.go

package errors

type LifeCycleError struct {
    Code    int
    Message string
    Context map[string]interface{}
}

func (e *LifeCycleError) Error() string {
    return fmt.Sprintf("LifeCycle error [%d]: %s", e.Code, e.Message)
}

const (
    ErrCodeInvalidState = iota + 1000
    ErrCodeEntityNotFound
    ErrCodeStateTransition
    ErrCodeSystemLocked
)

// 创建错误实例的辅助函数
func NewLifeCycleError(code int, message string) *LifeCycleError {
    return &LifeCycleError{
        Code:    code,
        Message: message,
        Context: make(map[string]interface{}),
    }
}
