package glogi

import (
	"fmt"
	"os"
)

// Backward compatibility with standard log package.
// These functions log at INFO level (or ERROR for Fatal/Panic).

// Print logs arguments at INFO level (like fmt.Print)
func Print(v ...any) {
	Info(fmt.Sprint(v...))
}

// Println logs arguments at INFO level (like fmt.Println)
func Println(v ...any) {
	Info(fmt.Sprint(v...))
}

// Printf logs formatted message at INFO level
func Printf(format string, v ...any) {
	Info(fmt.Sprintf(format, v...))
}

// Fatalln logs at FATAL level and exits
func Fatalln(v ...any) {
	ensureInit()
	logger.Log(nil, LevelFatal, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf logs formatted message at FATAL level and exits
func Fatalf(format string, v ...any) {
	ensureInit()
	logger.Log(nil, LevelFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Panic logs at PANIC level and panics (standard log.Panic signature)
func Panic(v ...any) {
	msg := fmt.Sprint(v...)
	ensureInit()
	logger.Log(nil, LevelPanic, msg)
	panic(msg)
}

// Panicln logs at PANIC level and panics
func Panicln(v ...any) {
	msg := fmt.Sprint(v...)
	ensureInit()
	logger.Log(nil, LevelPanic, msg)
	panic(msg)
}

// Panicf logs formatted message at PANIC level and panics
func Panicf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	ensureInit()
	logger.Log(nil, LevelPanic, msg)
	panic(msg)
}
