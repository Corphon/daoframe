//security/types.go
package security

import (
    "context"
    "time"
)

// SecurityLevel 安全级别
type SecurityLevel int

const (
    LevelLow SecurityLevel = iota
    LevelMedium
    LevelHigh
    LevelCritical
)

// Principal 身份主体
type Principal struct {
    ID        string
    Type      string
    Name      string
    Roles     []string
    Groups    []string
    Metadata  map[string]interface{}
    Created   time.Time
    LastLogin time.Time
}

// Permission 权限定义
type Permission struct {
    ID          string
    Resource    string
    Action      string
    Effect      EffectType
    Conditions  []Condition
    Priority    int
    ExpireAt    time.Time
}

type EffectType string

const (
    Allow EffectType = "allow"
    Deny  EffectType = "deny"
)

// Condition 条件接口
type Condition interface {
    Evaluate(ctx context.Context, principal *Principal) bool
}
