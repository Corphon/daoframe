// tools/logger/logger.go
package logger

// LoggerOption 日志选项
type LoggerOption struct {
    Level      LogLevel
    Format     LogFormat
    Output     []io.Writer
    TimeFormat string
    Buffer     int
    Async      bool
}

// DaoLogger 改进的日志器
type DaoLogger struct {
    opts     *LoggerOption
    handlers []LogHandler
    filters  []LogFilter
    buffer   *ring.Buffer
    manager  *LogManager
    metrics  *LogMetrics
    mu       sync.RWMutex
}

// LogManager 日志管理器
type LogManager struct {
    loggers   map[string]*DaoLogger
    rotation  *LogRotation
    writer    *LogWriter
    formatter LogFormatter
}

// LogRotation 日志轮转
type LogRotation struct {
    MaxSize    int64
    MaxBackups int
    MaxAge     int
    Compress   bool
}
