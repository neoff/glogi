package glogi

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// Default ANSI color codes
const (
	defaultColorReset     = "\033[0m"
	defaultColorDarkGray  = "\033[90m" // TRACE - dark gray
	defaultColorCyan      = "\033[36m" // DEBUG - cyan (distinguishable from white INFO)
	defaultColorLightGray = "\033[37m" // Old DEBUG
	defaultColorYellow    = "\033[33m" // WARN
	defaultColorRed       = "\033[31m" // ERROR, FATAL, PANIC
	defaultColorGreen     = "\033[32m" // Source location
)

// Configurable settings (can be overridden via env or SetXxx functions)
var (
	sourceWidth    = 20 // Default source width, configurable via LOG_SOURCE_WIDTH
	colorReset     = defaultColorReset
	colorTrace     = defaultColorDarkGray
	colorDebug     = defaultColorCyan
	colorInfo      = "" // No color for INFO
	colorWarn      = defaultColorYellow
	colorError     = defaultColorRed
	colorSource    = defaultColorGreen
	colorsDisabled = false
	configLoaded   = false
)

// initConfig reads configuration from environment variables
func initConfig() {
	if configLoaded {
		return
	}
	configLoaded = true

	// Source width
	if w := os.Getenv("LOG_SOURCE_WIDTH"); w != "" {
		if width, err := strconv.Atoi(w); err == nil && width > 0 {
			sourceWidth = width
		}
	}

	// Disable colors
	if os.Getenv("LOG_NO_COLOR") == "1" || os.Getenv("LOG_NO_COLOR") == "true" {
		colorsDisabled = true
	}

	// Custom colors (ANSI codes like "32" for green, or named colors)
	if c := os.Getenv("LOG_COLOR_TRACE"); c != "" {
		colorTrace = parseColor(c)
	}
	if c := os.Getenv("LOG_COLOR_DEBUG"); c != "" {
		colorDebug = parseColor(c)
	}
	if c := os.Getenv("LOG_COLOR_INFO"); c != "" {
		colorInfo = parseColor(c)
	}
	if c := os.Getenv("LOG_COLOR_WARN"); c != "" {
		colorWarn = parseColor(c)
	}
	if c := os.Getenv("LOG_COLOR_ERROR"); c != "" {
		colorError = parseColor(c)
	}
	if c := os.Getenv("LOG_COLOR_SOURCE"); c != "" {
		colorSource = parseColor(c)
	}
}

// parseColor converts color config to ANSI code
// Accepts: "32" (just code) or "\033[32m" (full ANSI) or "green" (named)
func parseColor(c string) string {
	c = strings.TrimSpace(c)
	if c == "" {
		return ""
	}
	// Named colors
	switch strings.ToLower(c) {
	case "red":
		return "\033[31m"
	case "green":
		return "\033[32m"
	case "yellow":
		return "\033[33m"
	case "blue":
		return "\033[34m"
	case "magenta":
		return "\033[35m"
	case "cyan":
		return "\033[36m"
	case "white":
		return "\033[37m"
	case "gray", "grey":
		return "\033[90m"
	case "none", "off":
		return ""
	}
	// If already contains escape sequence
	if strings.Contains(c, "\033") || strings.Contains(c, "\\033") {
		return strings.ReplaceAll(c, "\\033", "\033")
	}
	// Just a number - wrap in ANSI
	if _, err := strconv.Atoi(c); err == nil {
		return fmt.Sprintf("\033[%sm", c)
	}
	return c
}

// SetSourceWidth sets the fixed width for source location display
func SetSourceWidth(width int) {
	if width > 0 {
		sourceWidth = width
	}
}

// SetColorTrace sets the color for TRACE level
func SetColorTrace(color string) { colorTrace = parseColor(color) }

// SetColorDebug sets the color for DEBUG level
func SetColorDebug(color string) { colorDebug = parseColor(color) }

// SetColorInfo sets the color for INFO level
func SetColorInfo(color string) { colorInfo = parseColor(color) }

// SetColorWarn sets the color for WARN level
func SetColorWarn(color string) { colorWarn = parseColor(color) }

// SetColorError sets the color for ERROR level
func SetColorError(color string) { colorError = parseColor(color) }

// SetColorSource sets the color for source location
func SetColorSource(color string) { colorSource = parseColor(color) }

// DisableColors disables all color output
func DisableColors() { colorsDisabled = true }

// EnableColors enables color output
func EnableColors() { colorsDisabled = false }

// ColoredHandler implements slog.Handler with colored level output
type ColoredHandler struct {
	level  *slog.LevelVar
	writer io.Writer
	attrs  []slog.Attr
	groups []string
}

// NewColoredHandler creates a new colored handler
func NewColoredHandler(w io.Writer, level *slog.LevelVar) *ColoredHandler {
	initConfig() // Read config from env on first handler creation
	return &ColoredHandler{
		level:  level,
		writer: w,
	}
}

func (h *ColoredHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.level.Level()
}

func (h *ColoredHandler) Handle(_ context.Context, r slog.Record) error {
	// Format: [2025/12/26 15:04:05] LEVEL [source_location] message key=value...
	timeStr := r.Time.Format("2006/01/02 15:04:05")
	levelStr, levelColor := h.formatLevelWithColor(r.Level)

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
			loc := fmt.Sprintf("%s:%d", file, f.Line)
			// Pad or truncate to fixed width
			if len(loc) > sourceWidth {
				loc = loc[:sourceWidth]
			} else {
				loc = fmt.Sprintf("%-*s", sourceWidth, loc)
			}
			if !colorsDisabled && colorSource != "" {
				source = fmt.Sprintf("%s[%s]%s", colorSource, loc, colorReset)
			} else {
				source = fmt.Sprintf("[%s]", loc)
			}
		}
	}

	// Build message content (will be colorized)
	msgContent := r.Message

	// Add attributes
	r.Attrs(func(a slog.Attr) bool {
		msgContent += fmt.Sprintf(" %s=%v", a.Key, a.Value.Any())
		return true
	})

	// Add handler-level attrs
	for _, a := range h.attrs {
		msgContent += fmt.Sprintf(" %s=%v", a.Key, a.Value.Any())
	}

	// Apply level color to message content
	if !colorsDisabled && levelColor != "" {
		msgContent = fmt.Sprintf("%s%s%s", levelColor, msgContent, colorReset)
	}

	// Build final message: [time] LEVEL [source] message
	msg := fmt.Sprintf("[%s] %s %s %s\n", timeStr, levelStr, source, msgContent)

	_, err := h.writer.Write([]byte(msg))
	return err
}

func (h *ColoredHandler) formatLevelWithColor(l slog.Level) (string, string) {
	var name string
	var color string

	switch {
	case l <= LevelTrace:
		name = "TRACE"
		color = colorTrace
	case l <= LevelDebug:
		name = "DEBUG"
		color = colorDebug
	case l <= LevelInfo:
		name = "INFO"
		color = colorInfo
	case l <= LevelWarn:
		name = "WARN"
		color = colorWarn
	case l <= LevelError:
		name = "ERROR"
		color = colorError
	case l <= LevelFatal:
		name = "FATAL"
		color = colorError
	default:
		name = "PANIC"
		color = colorError
	}

	// Fixed width: 5 characters
	paddedName := fmt.Sprintf("%-5s", name)

	if colorsDisabled || color == "" {
		return paddedName, ""
	}
	return fmt.Sprintf("%s%s%s", color, paddedName, colorReset), color
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
