//config/loader.go
package config

import (
    "encoding/json"
    "io/ioutil"
    "path/filepath"
)

// ConfigLoader 配置加载器
type ConfigLoader struct {
    configDir   string
    fileFormat  string
    manager     *ConfigManager
}

// LoadOption 加载选项
type LoadOption struct {
    Required      bool
    DefaultConfig interface{}
    Validator     ConfigValidator
}

func NewConfigLoader(configDir string) *ConfigLoader {
    return &ConfigLoader{
        configDir:  configDir,
        fileFormat: "json",
        manager:    NewConfigManager(),
    }
}

// LoadConfig 加载指定类型的配置
func (cl *ConfigLoader) LoadConfig(typ ConfigType, dst interface{}, opt *LoadOption) error {
    filename := filepath.Join(cl.configDir, string(typ)+"."+cl.fileFormat)
    
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        if opt != nil && !opt.Required {
            // 使用默认配置
            if opt.DefaultConfig != nil {
                return cl.manager.RegisterConfig(typ, opt.DefaultConfig, opt.Validator)
            }
        }
        return err
    }
    
    if err := json.Unmarshal(data, dst); err != nil {
        return err
    }
    
    return cl.manager.RegisterConfig(typ, dst, opt.Validator)
}
