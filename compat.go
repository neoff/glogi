package glogi

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"
)

// Backward compatibility with standard log package.
// These functions log at INFO level (or ERROR for Fatal/Panic).

// logCompatWithCaller logs compat messages with correct caller
func logCompatWithCaller(lvl slog.Level, msg string) {
	ensureInit()
	if !logger.Enabled(context.Background(), lvl) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip: Callers, logCompatWithCaller, Print*/Fatal*/Panic*

	r := slog.NewRecord(time.Now(), lvl, msg, pcs[0])
	_ = logger.Handler().Handle(context.Background(), r)
}

// Print logs arguments at INFO level (like fmt.Print)
func Print(v ...any) {
	logCompatWithCaller(LevelInfo, fmt.Sprint(v...))
}

// Println logs arguments at INFO level (like fmt.Println)
func Println(v ...any) {
	logCompatWithCaller(LevelInfo, fmt.Sprint(v...))
}

// Printf logs formatted message at INFO level
func Printf(format string, v ...any) {
	logCompatWithCaller(LevelInfo, fmt.Sprintf(format, v...))
}

// Fatalln logs at FATAL level and exits
func Fatalln(v ...any) {
	logCompatWithCaller(LevelFatal, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf logs formatted message at FATAL level and exits
func Fatalf(format string, v ...any) {
	logCompatWithCaller(LevelFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Panic logs at PANIC level and panics (standard log.Panic signature)
func Panic(v ...any) {
	msg := fmt.Sprint(v...)
	logCompatWithCaller(LevelPanic, msg)
	panic(msg)
}

// Panicln logs at PANIC level and panics
func Panicln(v ...any) {
	msg := fmt.Sprint(v...)
	logCompatWithCaller(LevelPanic, msg)
	panic(msg)
}

// Panicf logs formatted message at PANIC level and panics
func Panicf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logCompatWithCaller(LevelPanic, msg)
	panic(msg)
}
