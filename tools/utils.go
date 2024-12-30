// tools/utils.go

package tools

import (
    "crypto/md5"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io"
    "math/rand"
    "os"
    "path/filepath"
    "runtime"
    "strings"
    "sync"
    "time"
)

// TimeFormat 常用时间格式
const (
    TimeFormatFull     = "2006-01-02 15:04:05.000"
    TimeFormatStandard = "2006-01-02 15:04:05"
    TimeFormatDate     = "2006-01-02"
    TimeFormatTime     = "15:04:05"
)

// ByteSize 字节大小单位
const (
    KB = 1024
    MB = 1024 * KB
    GB = 1024 * MB
    TB = 1024 * GB
)

// PathExists 检查路径是否存在
func PathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

// EnsureDir 确保目录存在，如果不存在则创建
func EnsureDir(path string) error {
    exists, err := PathExists(path)
    if err != nil {
        return err
    }
    if !exists {
        return os.MkdirAll(path, 0755)
    }
    return nil
}

// MD5 计算字符串的 MD5 值
func MD5(str string) string {
    h := md5.New()
    h.Write([]byte(str))
    return hex.EncodeToString(h.Sum(nil))
}

// SHA256 计算字符串的 SHA256 值
func SHA256(str string) string {
    h := sha256.New()
    h.Write([]byte(str))
    return hex.EncodeToString(h.Sum(nil))
}

// RandomString 生成指定长度的随机字符串
func RandomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[seededRand.Intn(len(charset))]
    }
    return string(b)
}

// FormatFileSize 格式化文件大小
func FormatFileSize(size int64) string {
    switch {
    case size >= TB:
        return fmt.Sprintf("%.2fTB", float64(size)/float64(TB))
    case size >= GB:
        return fmt.Sprintf("%.2fGB", float64(size)/float64(GB))
    case size >= MB:
        return fmt.Sprintf("%.2fMB", float64(size)/float64(MB))
    case size >= KB:
        return fmt.Sprintf("%.2fKB", float64(size)/float64(KB))
    default:
        return fmt.Sprintf("%dB", size)
    }
}

// GetCurrentDirectory 获取当前工作目录
func GetCurrentDirectory() string {
    dir, err := os.Getwd()
    if err != nil {
        return ""
    }
    return dir
}

// GetExecutablePath 获取可执行文件路径
func GetExecutablePath() (string, error) {
    exe, err := os.Executable()
    if err != nil {
        return "", err
    }
    return filepath.EvalSymlinks(exe)
}

// MemoryCache 简单的内存缓存实现
type MemoryCache struct {
    sync.RWMutex
    data    map[string]interface{}
    expires map[string]time.Time
}

// NewMemoryCache 创建新的内存缓存
func NewMemoryCache() *MemoryCache {
    cache := &MemoryCache{
        data:    make(map[string]interface{}),
        expires: make(map[string]time.Time),
    }
    go cache.cleanup()
    return cache
}

// Set 设置缓存
func (c *MemoryCache) Set(key string, value interface{}, duration time.Duration) {
    c.Lock()
    defer c.Unlock()
    
    c.data[key] = value
    if duration > 0 {
        c.expires[key] = time.Now().Add(duration)
    }
}

// Get 获取缓存
func (c *MemoryCache) Get(key string) (interface{}, bool) {
    c.RLock()
    defer c.RUnlock()
    
    if expire, ok := c.expires[key]; ok && time.Now().After(expire) {
        return nil, false
    }
    
    value, ok := c.data[key]
    return value, ok
}

// Delete 删除缓存
func (c *MemoryCache) Delete(key string) {
    c.Lock()
    defer c.Unlock()
    
    delete(c.data, key)
    delete(c.expires, key)
}

// cleanup 清理过期缓存
func (c *MemoryCache) cleanup() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        c.Lock()
        now := time.Now()
        for key, expire := range c.expires {
            if now.After(expire) {
                delete(c.data, key)
                delete(c.expires, key)
            }
        }
        c.Unlock()
    }
}

// AsyncRunner 异步任务运行器
type AsyncRunner struct {
    wg sync.WaitGroup
}

// NewAsyncRunner 创建异步任务运行器
func NewAsyncRunner() *AsyncRunner {
    return &AsyncRunner{}
}

// Run 运行异步任务
func (r *AsyncRunner) Run(fn func()) {
    r.wg.Add(1)
    go func() {
        defer r.wg.Done()
        fn()
    }()
}

// Wait 等待所有任务完成
func (r *AsyncRunner) Wait() {
    r.wg.Wait()
}

// RuntimeStats 运行时统计信息
type RuntimeStats struct {
    GoVersion    string
    GOOS         string
    GOARCH       string
    NumCPU       int
    NumGoroutine int
    MemStats     runtime.MemStats
}

// GetRuntimeStats 获取运行时统计信息
func GetRuntimeStats() RuntimeStats {
    var stats RuntimeStats
    var mem runtime.MemStats
    
    runtime.ReadMemStats(&mem)
    
    stats.GoVersion = runtime.Version()
    stats.GOOS = runtime.GOOS
    stats.GOARCH = runtime.GOARCH
    stats.NumCPU = runtime.NumCPU()
    stats.NumGoroutine = runtime.NumGoroutine()
    stats.MemStats = mem
    
    return stats
}

// Retry 重试执行函数
func Retry(attempts int, sleep time.Duration, fn func() error) error {
    var err error
    
    for i := 0; i < attempts; i++ {
        if err = fn(); err == nil {
            return nil
        }
        
        if i < attempts-1 {
            time.Sleep(sleep)
        }
    }
    
    return fmt.Errorf("在%d次尝试后失败: %w", attempts, err)
}

// DeepCopy 深拷贝
func DeepCopy(src interface{}, dst interface{}) error {
    if data, err := json.Marshal(src); err != nil {
        return err
    } else {
        return json.Unmarshal(data, dst)
    }
}

// TruncateString 截断字符串
func TruncateString(str string, length int) string {
    if length <= 0 {
        return ""
    }
    
    if len(str) <= length {
        return str
    }
    
    return str[:length] + "..."
}

// IsValidEmail 验证邮箱格式
func IsValidEmail(email string) bool {
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false
    }
    
    if len(parts[0]) == 0 || len(parts[1]) == 0 {
        return false
    }
    
    if !strings.Contains(parts[1], ".") {
        return false
    }
    
    return true
}
