// tools/logger.go

package tools

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "runtime"
    "sync"
    "time"
)

// LogLevel 日志级别
type LogLevel int

const (
    DebugLevel LogLevel = iota
    InfoLevel
    WarnLevel
    ErrorLevel
    FatalLevel
)

var levelNames = map[LogLevel]string{
    DebugLevel: "DEBUG",
    InfoLevel:  "INFO",
    WarnLevel:  "WARN",
    ErrorLevel: "ERROR",
    FatalLevel: "FATAL",
}

// DaoLogger 道家风格的日志器
type DaoLogger struct {
    mu       sync.Mutex
    logger   *log.Logger
    level    LogLevel
    outputs  []io.Writer
    filename string
}

// LoggerOption 日志配置选项
type LoggerOption func(*DaoLogger)

// NewDaoLogger 创建新的日志器
func NewDaoLogger(options ...LoggerOption) *DaoLogger {
    dl := &DaoLogger{
        level:   InfoLevel,
        outputs: []io.Writer{os.Stdout},
    }

    // 应用配置选项
    for _, option := range options {
        option(dl)
    }

    // 创建多输出写入器
    multiWriter := io.MultiWriter(dl.outputs...)
    dl.logger = log.New(multiWriter, "", 0)

    return dl
}

// WithLevel 设置日志级别
func WithLevel(level LogLevel) LoggerOption {
    return func(dl *DaoLogger) {
        dl.level = level
    }
}

// WithFile 添加文件输出
func WithFile(filename string) LoggerOption {
    return func(dl *DaoLogger) {
        if filename != "" {
            // 确保日志目录存在
            dir := filepath.Dir(filename)
            if err := os.MkdirAll(dir, 0755); err == nil {
                if file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644); err == nil {
                    dl.outputs = append(dl.outputs, file)
                    dl.filename = filename
                }
            }
        }
    }
}

// log 内部日志方法
func (dl *DaoLogger) log(level LogLevel, format string, args ...interface{}) {
    if level < dl.level {
        return
    }

    dl.mu.Lock()
    defer dl.mu.Unlock()

    // 获取调用信息
    _, file, line, ok := runtime.Caller(2)
    if !ok {
        file = "unknown"
        line = 0
    }

    // 格式化时间
    now := time.Now().Format("2006-01-02 15:04:05.000")
    
    // 构建日志前缀
    prefix := fmt.Sprintf("[%s][%s][%s:%d] ", 
        now,
        levelNames[level],
        filepath.Base(file),
        line,
    )

    // 格式化消息
    var msg string
    if len(args) > 0 {
        msg = fmt.Sprintf(format, args...)
    } else {
        msg = format
    }

    // 写入日志
    dl.logger.Printf("%s%s", prefix, msg)
}

// Debug 调试日志
func (dl *DaoLogger) Debug(format string, args ...interface{}) {
    dl.log(DebugLevel, format, args...)
}

// Info 信息日志
func (dl *DaoLogger) Info(format string, args ...interface{}) {
    dl.log(InfoLevel, format, args...)
}

// Warn 警告日志
func (dl *DaoLogger) Warn(format string, args ...interface{}) {
    dl.log(WarnLevel, format, args...)
}

// Error 错误日志
func (dl *DaoLogger) Error(format string, args ...interface{}) {
    dl.log(ErrorLevel, format, args...)
}

// Fatal 致命错误日志
func (dl *DaoLogger) Fatal(format string, args ...interface{}) {
    dl.log(FatalLevel, format, args...)
    os.Exit(1)
}

// Rotate 日志轮转
func (dl *DaoLogger) Rotate() error {
    if dl.filename == "" {
        return nil
    }

    dl.mu.Lock()
    defer dl.mu.Unlock()

    // 关闭当前日志文件
    for _, output := range dl.outputs {
        if file, ok := output.(*os.File); ok && file.Name() == dl.filename {
            file.Close()
            break
        }
    }

    // 生成新的文件名
    timestamp := time.Now().Format("20060102150405")
    newFilename := fmt.Sprintf("%s.%s", dl.filename, timestamp)
    
    // 重命名当前日志文件
    if err := os.Rename(dl.filename, newFilename); err != nil {
        return err
    }

    // 创建新的日志文件
    newFile, err := os.OpenFile(dl.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    // 更新输出列表
    newOutputs := make([]io.Writer, 0, len(dl.outputs))
    for _, output := range dl.outputs {
        if file, ok := output.(*os.File); !ok || file.Name() != dl.filename {
            newOutputs = append(newOutputs, output)
        }
    }
    newOutputs = append(newOutputs, newFile)
    dl.outputs = newOutputs

    // 更新logger
    dl.logger.SetOutput(io.MultiWriter(dl.outputs...))

    return nil
}

// 全局默认日志器
var DefaultLogger = NewDaoLogger()

// 全局日志函数
func Debug(format string, args ...interface{}) {
    DefaultLogger.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
    DefaultLogger.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
    DefaultLogger.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
    DefaultLogger.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
    DefaultLogger.Fatal(format, args...)
}
