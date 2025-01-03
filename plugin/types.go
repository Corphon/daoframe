//plugin/types.go
package plugin

import (
    "context"
    "time"
    "github.com/Corphon/daoframe/errors"
)

// PluginState 插件状态
type PluginState int

const (
    PluginStateUnknown PluginState = iota
    PluginStateInitializing
    PluginStateActive
    PluginStatePaused
    PluginStateError
    PluginStateTerminated
)

// PluginInfo 插件信息
type PluginInfo struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Version     string            `json:"version"`
    Author      string            `json:"author"`
    Description string            `json:"description"`
    Homepage    string            `json:"homepage"`
    License     string            `json:"license"`
    Tags        []string          `json:"tags"`
    Metadata    map[string]string `json:"metadata"`
    Dependencies []Dependency     `json:"dependencies"`
    Created     time.Time         `json:"created"`
    Updated     time.Time         `json:"updated"`
}

// Dependency 插件依赖
type Dependency struct {
    ID          string `json:"id"`
    Version     string `json:"version"`
    Required    bool   `json:"required"`
    MinVersion  string `json:"min_version,omitempty"`
    MaxVersion  string `json:"max_version,omitempty"`
}
