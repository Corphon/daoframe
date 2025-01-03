// plugin/plugin.go
package plugin

type Plugin interface {
    ID() string
    Version() string
    Load() error
    Unload() error
}

// plugin/manager.go
type PluginManager struct {
    plugins     map[string]Plugin
    loader      *PluginLoader
    registry    *PluginRegistry
    hooks       map[string][]PluginHook
}

// plugin/loader.go
type PluginLoader struct {
    paths      []string
    validators []PluginValidator
    cache      *cache.Cache
}
