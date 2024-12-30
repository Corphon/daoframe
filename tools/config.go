// tools/config.go

package tools

import (
    "encoding/json"
    "encoding/xml"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
    "yaml"
)

// ConfigFormat 配置文件格式
type ConfigFormat string

const (
    FormatJSON ConfigFormat = "json"
    FormatXML  ConfigFormat = "xml"
    FormatYAML ConfigFormat = "yaml"
)

// ConfigOption 配置选项
type ConfigOption struct {
    AutoReload     bool          // 自动重载
    ReloadInterval time.Duration // 重载间隔
    BackupCount    int          // 备份数量
    Format         ConfigFormat  // 配置格式
}

// DaoConfig 配置管理器
type DaoConfig struct {
    mu            sync.RWMutex
    data          map[string]interface{}
    filepath      string
    lastModified  time.Time
    option        ConfigOption
    watchers      []ConfigWatcher
    stopChan      chan struct{}
}

// ConfigWatcher 配置变更观察者
type ConfigWatcher interface {
    OnConfigChange(oldConfig, newConfig map[string]interface{})
}

// NewDaoConfig 创建新的配置管理器
func NewDaoConfig(filepath string, opt ConfigOption) (*DaoConfig, error) {
    dc := &DaoConfig{
        data:     make(map[string]interface{}),
        filepath: filepath,
        option:   opt,
        watchers: make([]ConfigWatcher, 0),
        stopChan: make(chan struct{}),
    }

    // 首次加载配置
    if err := dc.Load(); err != nil {
        return nil, err
    }

    // 如果启用自动重载，启动监控
    if opt.AutoReload {
        go dc.watchConfig()
    }

    return dc, nil
}

// Load 加载配置文件
func (dc *DaoConfig) Load() error {
    dc.mu.Lock()
    defer dc.mu.Unlock()

    file, err := os.Open(dc.filepath)
    if err != nil {
        return fmt.Errorf("打开配置文件失败: %w", err)
    }
    defer file.Close()

    // 读取文件内容
    content, err := io.ReadAll(file)
    if err != nil {
        return fmt.Errorf("读取配置文件失败: %w", err)
    }

    // 备份当前配置
    oldConfig := make(map[string]interface{})
    for k, v := range dc.data {
        oldConfig[k] = v
    }

    // 根据格式解析配置
    newData := make(map[string]interface{})
    switch dc.option.Format {
    case FormatJSON:
        err = json.Unmarshal(content, &newData)
    case FormatXML:
        err = xml.Unmarshal(content, &newData)
    case FormatYAML:
        err = yaml.Unmarshal(content, &newData)
    default:
        return fmt.Errorf("不支持的配置格式: %s", dc.option.Format)
    }

    if err != nil {
        return fmt.Errorf("解析配置文件失败: %w", err)
    }

    // 更新配置
    dc.data = newData
    dc.lastModified = time.Now()

    // 通知观察者
    dc.notifyWatchers(oldConfig, newData)

    return nil
}

// Save 保存配置到文件
func (dc *DaoConfig) Save() error {
    dc.mu.RLock()
    defer dc.mu.RUnlock()

    // 如果需要备份，创建备份文件
    if dc.option.BackupCount > 0 {
        if err := dc.createBackup(); err != nil {
            return err
        }
    }

    // 根据格式编码配置
    var content []byte
    var err error

    switch dc.option.Format {
    case FormatJSON:
        content, err = json.MarshalIndent(dc.data, "", "  ")
    case FormatXML:
        content, err = xml.MarshalIndent(dc.data, "", "  ")
    case FormatYAML:
        content, err = yaml.Marshal(dc.data)
    default:
        return fmt.Errorf("不支持的配置格式: %s", dc.option.Format)
    }

    if err != nil {
        return fmt.Errorf("编码配置失败: %w", err)
    }

    // 写入文件
    if err := os.WriteFile(dc.filepath, content, 0644); err != nil {
        return fmt.Errorf("写入配置文件失败: %w", err)
    }

    return nil
}

// Get 获取配置项
func (dc *DaoConfig) Get(key string) (interface{}, bool) {
    dc.mu.RLock()
    defer dc.mu.RUnlock()

    value, exists := dc.data[key]
    return value, exists
}

// Set 设置配置项
func (dc *DaoConfig) Set(key string, value interface{}) {
    dc.mu.Lock()
    defer dc.mu.Unlock()

    oldConfig := make(map[string]interface{})
    for k, v := range dc.data {
        oldConfig[k] = v
    }

    dc.data[key] = value
    dc.notifyWatchers(oldConfig, dc.data)
}

// Delete 删除配置项
func (dc *DaoConfig) Delete(key string) {
    dc.mu.Lock()
    defer dc.mu.Unlock()

    oldConfig := make(map[string]interface{})
    for k, v := range dc.data {
        oldConfig[k] = v
    }

    delete(dc.data, key)
    dc.notifyWatchers(oldConfig, dc.data)
}

// AddWatcher 添加配置观察者
func (dc *DaoConfig) AddWatcher(watcher ConfigWatcher) {
    dc.mu.Lock()
    defer dc.mu.Unlock()
    dc.watchers = append(dc.watchers, watcher)
}

// notifyWatchers 通知所有观察者
func (dc *DaoConfig) notifyWatchers(oldConfig, newConfig map[string]interface{}) {
    for _, watcher := range dc.watchers {
        go watcher.OnConfigChange(oldConfig, newConfig)
    }
}

// watchConfig 监控配置文件变化
func (dc *DaoConfig) watchConfig() {
    ticker := time.NewTicker(dc.option.ReloadInterval)
    defer ticker.Stop()

    for {
        select {
        case <-dc.stopChan:
            return
        case <-ticker.C:
            if dc.isFileModified() {
                if err := dc.Load(); err != nil {
                    DefaultLogger.Error("重载配置失败: %v", err)
                }
            }
        }
    }
}

// isFileModified 检查文件是否被修改
func (dc *DaoConfig) isFileModified() bool {
    info, err := os.Stat(dc.filepath)
    if err != nil {
        return false
    }
    return info.ModTime().After(dc.lastModified)
}

// createBackup 创建配置文件备份
func (dc *DaoConfig) createBackup() error {
    // 生成备份文件名
    timestamp := time.Now().Format("20060102150405")
    backupPath := fmt.Sprintf("%s.%s", dc.filepath, timestamp)

    // 创建备份文件
    if err := copyFile(dc.filepath, backupPath); err != nil {
        return fmt.Errorf("创建配置备份失败: %w", err)
    }

    // 清理旧备份
    return dc.cleanOldBackups()
}

// cleanOldBackups 清理旧的备份文件
func (dc *DaoConfig) cleanOldBackups() error {
    dir := filepath.Dir(dc.filepath)
    base := filepath.Base(dc.filepath)

    // 获取所有备份文件
    files, err := filepath.Glob(filepath.Join(dir, base+".*"))
    if err != nil {
        return err
    }

    // 如果备份数量超过限制，删除最旧的文件
    if len(files) > dc.option.BackupCount {
        // 按修改时间排序
        type backup struct {
            path    string
            modTime time.Time
        }
        backups := make([]backup, 0, len(files))

        for _, file := range files {
            info, err := os.Stat(file)
            if err != nil {
                continue
            }
            backups = append(backups, backup{file, info.ModTime()})
        }

        // 排序
        sort := func(i, j int) bool {
            return backups[i].modTime.Before(backups[j].modTime)
        }
        for i := 0; i < len(backups)-1; i++ {
            for j := i + 1; j < len(backups); j++ {
                if !sort(i, j) {
                    backups[i], backups[j] = backups[j], backups[i]
                }
            }
        }

        // 删除多余的备份
        for i := 0; i < len(backups)-dc.option.BackupCount; i++ {
            if err := os.Remove(backups[i].path); err != nil {
                DefaultLogger.Error("删除旧备份失败: %v", err)
            }
        }
    }

    return nil
}

// copyFile 复制文件
func copyFile(src, dst string) error {
    source, err := os.Open(src)
    if err != nil {
        return err
    }
    defer source.Close()

    destination, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destination.Close()

    _, err = io.Copy(destination, source)
    return err
}

// Close 关闭配置管理器
func (dc *DaoConfig) Close() error {
    close(dc.stopChan)
    return dc.Save()
}
