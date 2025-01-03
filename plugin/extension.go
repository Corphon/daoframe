//plugin/extension.go
package plugin

// Extension 扩展接口
type Extension interface {
    ID() string
    Type() string
    Version() string
    Init(ctx context.Context) error
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
}

// ExtensionPoint 扩展点接口
type ExtensionPoint interface {
    ID() string
    Description() string
    AddExtension(ext Extension) error
    RemoveExtension(id string) error
    GetExtensions() []Extension
}

// ExtensionRegistry 扩展注册表
type ExtensionRegistry struct {
    mu sync.RWMutex
    extensions map[string][]Extension
    points    map[string]ExtensionPoint
}
