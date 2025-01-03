//errors/error_codes.go
package errors

// System Error Codes (1000-1999)
const (
    ErrSystemInit        ErrorCode = 1000 + iota
    ErrSystemShutdown
    ErrSystemOverload
    ErrSystemTimeout
    ErrSystemResource
    ErrSystemState
    ErrSystemConfig
)

// Core Error Codes (2000-2999)
const (
    ErrCoreInit         ErrorCode = 2000 + iota
    ErrCoreState
    ErrCoreContext
    ErrCoreAdapt
    ErrCoreUniverse
)

// Model Error Codes (3000-3999)
const (
    ErrModelInit        ErrorCode = 3000 + iota
    ErrModelValidation
    ErrModelState
    ErrModelTransform
    ErrModelInteraction
)

// Lifecycle Error Codes (4000-4999)
const (
    ErrLifecycleInit    ErrorCode = 4000 + iota
    ErrLifecycleState
    ErrLifecycleEntity
    ErrLifecycleTransition
    ErrLifecycleLock
)

// Runtime Error Codes (5000-5999)
const (
    ErrRuntimeExec      ErrorCode = 5000 + iota
    ErrRuntimeResource
    ErrRuntimeTimeout
    ErrRuntimeState
    ErrRuntimeIO
)
