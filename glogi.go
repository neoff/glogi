// Package glogi provides structured logging with colored output and level-based filtering.
// It is designed as a drop-in replacement for the standard log package.
//
// Usage:
//
//	import log "github.com/neoff/glogi"
//
//	func main() {
//	    log.Init() // Reads LOG_LEVEL from env (default: INFO)
//	    log.Info("server started", "port", 8080)
//	}
package glogi

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"sync"
)

var (
	logger   *slog.Logger
	level    *slog.LevelVar
	initOnce sync.Once
	isInit   bool
)

// Custom log levels
const (
	LevelTrace = slog.Level(-8)
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
	LevelFatal = slog.Level(12)
	LevelPanic = slog.Level(16)
)

// Init initializes the global logger.
// Reads LOG_LEVEL from environment variable (default: INFO).
// Valid values: TRACE, DEBUG, INFO, WARN, ERROR
func Init() {
	initOnce.Do(func() {
		level = &slog.LevelVar{}
		level.Set(parseLevel(os.Getenv("LOG_LEVEL")))

		handler := NewColoredHandler(os.Stdout, level)
		logger = slog.New(handler)
		slog.SetDefault(logger)
		isInit = true
	})
}

// SetLevel changes the minimum log level at runtime
func SetLevel(l string) {
	if level != nil {
		level.Set(parseLevel(l))
	}
}

func parseLevel(s string) slog.Level {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "TRACE":
		return LevelTrace
	case "DEBUG":
		return LevelDebug
	case "WARN", "WARNING":
		return LevelWarn
	case "ERROR":
		return LevelError
	default:
		return LevelInfo
	}
}

func ensureInit() {
	if !isInit {
		Init()
	}
}

// Trace logs at TRACE level (light gray)
func Trace(msg string, args ...any) {
	ensureInit()
	logger.Log(context.Background(), LevelTrace, msg, args...)
}

// Debug logs at DEBUG level (gray)
func Debug(msg string, args ...any) {
	ensureInit()
	logger.Debug(msg, args...)
}

// Info logs at INFO level (no color)
func Info(msg string, args ...any) {
	ensureInit()
	logger.Info(msg, args...)
}

// Warn logs at WARN level (yellow)
func Warn(msg string, args ...any) {
	ensureInit()
	logger.Warn(msg, args...)
}

// Error logs at ERROR level (red)
func Error(msg string, args ...any) {
	ensureInit()
	logger.Error(msg, args...)
}

// Fatal logs at FATAL level (red) and calls os.Exit(1)
func Fatal(msg string, args ...any) {
	ensureInit()
	logger.Log(context.Background(), LevelFatal, msg, args...)
	os.Exit(1)
}

// Panic logs at PANIC level (red) and panics
func PanicLog(msg string, args ...any) {
	ensureInit()
	logger.Log(context.Background(), LevelPanic, msg, args...)
	panic(msg)
}

// Recover catches panic and logs it with stack trace. Use in defer.
func Recover() {
	if r := recover(); r != nil {
		ensureInit()
		buf := make([]byte, 4096)
		n := runtime.Stack(buf, false)
		logger.Log(context.Background(), LevelPanic, fmt.Sprintf("recovered: %v", r), "stack", string(buf[:n]))
	}
}
