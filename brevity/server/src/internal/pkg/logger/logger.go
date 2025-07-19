package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Level represents the logging level
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

var levelNames = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	FatalLevel: "FATAL",
}

// Field represents a key-value pair in the log
type Field struct {
	Key   string
	Value interface{}
}

// Logger is the interface for the logger
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Sync() error
	With(fields ...Field) Logger
}

// Config holds the configuration for the logger
type Config struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	FilePath string `mapstructure:"file_path"`
}

// customLogger implements Logger interface
type customLogger struct {
	level    Level
	format   string
	writer   io.Writer
	mu       sync.Mutex
	fields   []Field
	file     *os.File
	useFile  bool
	colorful bool
}

var (
	instance Logger
	once     sync.Once
	initErr  error
)

// Init initializes the logger with the given configuration
func Init(cfg Config) error {
	once.Do(func() {
		// Convert string level to Level type
		var level Level
		switch strings.ToLower(cfg.Level) {
		case "debug":
			level = DebugLevel
		case "info":
			level = InfoLevel
		case "warn", "warning":
			level = WarnLevel
		case "error":
			level = ErrorLevel
		default:
			level = InfoLevel
		}

		var writer io.Writer = os.Stdout
		var file *os.File
		useFile := false

		if cfg.FilePath != "" {
			if err := os.MkdirAll(filepath.Dir(cfg.FilePath), 0755); err != nil {
				initErr = fmt.Errorf("failed to create log directory: %w", err)
				return
			}

			file, initErr = os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if initErr != nil {
				initErr = fmt.Errorf("failed to open log file: %w", initErr)
				return
			}

			if cfg.Format == "console" {
				writer = io.MultiWriter(os.Stdout, file)
			} else {
				writer = file
			}
			useFile = true
		}

		instance = &customLogger{
			level:    level,
			format:   cfg.Format,
			writer:   writer,
			file:     file,
			useFile:  useFile,
			colorful: cfg.Format == "console",
		}
	})

	return initErr
}

// Get returns the default logger instance
func Get() Logger {
	if instance == nil {
		// Initialize with default configuration if not initialized
		_ = Init(Config{
			Level:  "info",
			Format: "console",
		})
	}
	return instance
}

// Sync flushes any buffered log entries
func Sync() error {
	if instance != nil {
		return instance.Sync()
	}
	return nil
}

// Implement Logger interface methods
func (l *customLogger) Debug(msg string, fields ...Field) {
	l.log(DebugLevel, msg, fields...)
}

func (l *customLogger) Info(msg string, fields ...Field) {
	l.log(InfoLevel, msg, fields...)
}

func (l *customLogger) Warn(msg string, fields ...Field) {
	l.log(WarnLevel, msg, fields...)
}

func (l *customLogger) Error(msg string, fields ...Field) {
	l.log(ErrorLevel, msg, fields...)
}

func (l *customLogger) Fatal(msg string, fields ...Field) {
	l.log(FatalLevel, msg, fields...)
	os.Exit(1)
}

func (l *customLogger) Debugf(format string, args ...interface{}) {
	l.log(DebugLevel, fmt.Sprintf(format, args...))
}

func (l *customLogger) Infof(format string, args ...interface{}) {
	l.log(InfoLevel, fmt.Sprintf(format, args...))
}

func (l *customLogger) Warnf(format string, args ...interface{}) {
	l.log(WarnLevel, fmt.Sprintf(format, args...))
}

func (l *customLogger) Errorf(format string, args ...interface{}) {
	l.log(ErrorLevel, fmt.Sprintf(format, args...))
}

func (l *customLogger) Fatalf(format string, args ...interface{}) {
	l.log(FatalLevel, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (l *customLogger) Sync() error {
	if l.useFile && l.file != nil {
		return l.file.Sync()
	}
	return nil
}

func (l *customLogger) With(fields ...Field) Logger {
	return &customLogger{
		level:    l.level,
		format:   l.format,
		writer:   l.writer,
		fields:   append(l.fields, fields...),
		file:     l.file,
		useFile:  l.useFile,
		colorful: l.colorful,
	}
}

// log is the internal logging function
func (l *customLogger) log(level Level, msg string, fields ...Field) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().Format(time.RFC3339)
	levelName := levelNames[level]
	allFields := append(l.fields, fields...)

	var logLine string
	if l.format == "json" {
		logLine = l.formatJSON(now, levelName, msg, allFields)
	} else {
		logLine = l.formatConsole(now, levelName, msg, allFields)
	}

	fmt.Fprintln(l.writer, logLine)
}

func (l *customLogger) formatJSON(timestamp, level, msg string, fields []Field) string {
	entry := fmt.Sprintf(`{"time":"%s","level":"%s","msg":"%s"`, timestamp, level, msg)

	for _, field := range fields {
		entry += fmt.Sprintf(`,"%s":%v`, field.Key, formatValue(field.Value))
	}

	entry += "}"
	return entry
}

func (l *customLogger) formatConsole(timestamp, level, msg string, fields []Field) string {
	var colorStart, colorReset string

	if l.colorful {
		colorStart, colorReset = getLevelColor(level)
	}

	entry := fmt.Sprintf("%s%s %s %s%s", colorStart, timestamp, level, msg, colorReset)

	for _, field := range fields {
		entry += fmt.Sprintf(" %s=%v", field.Key, field.Value)
	}

	return entry
}

func getLevelColor(level string) (string, string) {
	switch level {
	case "DEBUG":
		return "\033[36m", "\033[0m" // Cyan
	case "INFO":
		return "\033[32m", "\033[0m" // Green
	case "WARN":
		return "\033[33m", "\033[0m" // Yellow
	case "ERROR", "FATAL":
		return "\033[31m", "\033[0m" // Red
	default:
		return "", ""
	}
}

func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf(`"%s"`, v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Field constructors
func String(key, val string) Field {
	return Field{Key: key, Value: val}
}

func Int(key string, val int) Field {
	return Field{Key: key, Value: val}
}

func Int64(key string, val int64) Field {
	return Field{Key: key, Value: val}
}

func Float64(key string, val float64) Field {
	return Field{Key: key, Value: val}
}

func Bool(key string, val bool) Field {
	return Field{Key: key, Value: val}
}

func Duration(key string, val time.Duration) Field {
	return Field{Key: key, Value: val}
}

func Time(key string, val time.Time) Field {
	return Field{Key: key, Value: val}
}

func Any(key string, val interface{}) Field {
	return Field{Key: key, Value: val}
}

func ErrorField(err error) Field {
	return Field{Key: "error", Value: err.Error()}
}

func NamedError(key string, err error) Field {
	return Field{Key: key, Value: err.Error()}
}
