//config/types.go
package config

import "time"

// ConfigType 配置类型
type ConfigType string

const (
    CoreConfig     ConfigType = "core"
    SystemConfig   ConfigType = "system"
    ModelConfig    ConfigType = "model"
    StorageConfig  ConfigType = "storage"
)

// Duration 自定义时间Duration类型，支持字符串解析
type Duration struct {
    time.Duration
}

// UnmarshalJSON 实现json反序列化
func (d *Duration) UnmarshalJSON(b []byte) error {
    var v interface{}
    if err := json.Unmarshal(b, &v); err != nil {
        return err
    }
    
    switch value := v.(type) {
    case float64:
        d.Duration = time.Duration(value)
        return nil
    case string:
        var err error
        d.Duration, err = time.ParseDuration(value)
        if err != nil {
            return err
        }
        return nil
    default:
        return errors.New("invalid duration")
    }
}
