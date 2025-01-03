//tools/common/types.go
package common

// Result 通用结果类型
type Result[T any] struct {
    Data    T
    Error   error
    Success bool
}

// Pair 通用键值对类型
type Pair[K, V any] struct {
    Key   K
    Value V
}

// Queue 线程安全的队列
type Queue[T any] struct {
    items []T
    mu    sync.RWMutex
}
