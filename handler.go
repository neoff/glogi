package glogi

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strings"
)

// ANSI color codes
const (
	colorReset     = "\033[0m"
	colorLightGray = "\033[37m" // TRACE - very light
	colorGray      = "\033[90m" // DEBUG
	colorYellow    = "\033[33m" // WARN
	colorRed       = "\033[31m" // ERROR, FATAL, PANIC
)

// ColoredHandler implements slog.Handler with colored level output
type ColoredHandler struct {
	level  *slog.LevelVar
	writer io.Writer
	attrs  []slog.Attr
	groups []string
}

// NewColoredHandler creates a new colored handler
func NewColoredHandler(w io.Writer, level *slog.LevelVar) *ColoredHandler {
	return &ColoredHandler{
		level:  level,
		writer: w,
	}
}

func (h *ColoredHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.level.Level()
}

func (h *ColoredHandler) Handle(_ context.Context, r slog.Record) error {
	// Format: 2025/12/26 15:04:05 main.go:190: [LEVEL] message key=value...
	timeStr := r.Time.Format("2006/01/02 15:04:05")
	levelStr := h.formatLevel(r.Level)

	// Get source location from PC
	source := ""
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.File != "" {
			// Extract just the filename, not full path
			file := f.File
			if idx := strings.LastIndex(file, "/"); idx >= 0 {
				file = file[idx+1:]
			}
			source = fmt.Sprintf("%s:%d: ", file, f.Line)
		}
	}

	// Build message
	msg := fmt.Sprintf("%s %s%s %s", timeStr, source, levelStr, r.Message)

	// Add attributes
	r.Attrs(func(a slog.Attr) bool {
		msg += fmt.Sprintf(" %s=%v", a.Key, a.Value.Any())
		return true
	})

	// Add handler-level attrs
	for _, a := range h.attrs {
		msg += fmt.Sprintf(" %s=%v", a.Key, a.Value.Any())
	}

	msg += "\n"

	_, err := h.writer.Write([]byte(msg))
	return err
}

func (h *ColoredHandler) formatLevel(l slog.Level) string {
	var name string
	var color string

	switch {
	case l <= LevelTrace:
		name = "TRACE"
		color = colorLightGray
	case l <= LevelDebug:
		name = "DEBUG"
		color = colorGray
	case l <= LevelInfo:
		name = "INFO"
		color = ""
	case l <= LevelWarn:
		name = "WARN"
		color = colorYellow
	case l <= LevelError:
		name = "ERROR"
		color = colorRed
	case l <= LevelFatal:
		name = "FATAL"
		color = colorRed
	default:
		name = "PANIC"
		color = colorRed
	}

	if color == "" {
		return fmt.Sprintf("[%s]", name)
	}
	return fmt.Sprintf("%s[%s]%s", color, name, colorReset)
}

func (h *ColoredHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ColoredHandler{
		level:  h.level,
		writer: h.writer,
		attrs:  append(h.attrs, attrs...),
		groups: h.groups,
	}
}

func (h *ColoredHandler) WithGroup(name string) slog.Handler {
	return &ColoredHandler{
		level:  h.level,
		writer: h.writer,
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}

// Ensure ColoredHandler implements slog.Handler
var _ slog.Handler = (*ColoredHandler)(nil)
