# glogi

Structured logging library for Go with colored output and level-based filtering.
Drop-in replacement for standard `log` package.

## Installation

```bash
go get github.com/neoff/glogi@latest
```

## Usage

```go
import log "github.com/neoff/glogi"

func main() {
    defer log.Recover() // Catch panics
    log.Init()          // Reads LOG_LEVEL and other config from env
    
    // Structured logging (new style)
    log.Trace("trace message", "key", "value")
    log.Debug("debug info", "user_id", 123)
    log.Info("server started", "port", 8080)
    log.Warn("slow query", "duration", "500ms")
    log.Error("connection failed", "err", err)
    
    // Standard log compatibility (old style)
    log.Println("Hello, World!")
    log.Printf("Port: %d", 8080)
    log.Fatalf("Cannot start: %v", err)
}
```

## Log Levels

| Level | Color | Description |
|-------|-------|-------------|
| TRACE | Light gray | Detailed tracing |
| DEBUG | Gray | Debug information |
| INFO | No color | General information |
| WARN | Yellow | Warnings |
| ERROR | Red | Errors |
| FATAL | Red | Fatal + os.Exit(1) |
| PANIC | Red | Panic + stack trace |

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | `INFO` | Minimum log level (TRACE, DEBUG, INFO, WARN, ERROR) |
| `LOG_SOURCE_WIDTH` | `20` | Fixed width for source location column |
| `LOG_NO_COLOR` | `false` | Disable all colors (`1` or `true`) |
| `LOG_COLOR_TRACE` | `white` | Color for TRACE level |
| `LOG_COLOR_DEBUG` | `gray` | Color for DEBUG level |
| `LOG_COLOR_INFO` | (none) | Color for INFO level |
| `LOG_COLOR_WARN` | `yellow` | Color for WARN level |
| `LOG_COLOR_ERROR` | `red` | Color for ERROR level |
| `LOG_COLOR_SOURCE` | `green` | Color for source location |

### Color Values

Colors can be specified as:
- Named colors: `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`, `gray`
- ANSI code number: `32` (green), `31` (red), etc.
- Full ANSI sequence: `\033[32m`
- `none` or `off` to disable

### Example

```bash
# Development with wide source column
LOG_LEVEL=DEBUG LOG_SOURCE_WIDTH=25 ./myapp

# Production with no colors (for log aggregators)
LOG_LEVEL=INFO LOG_NO_COLOR=1 ./myapp

# Custom colors
LOG_COLOR_SOURCE=cyan LOG_COLOR_WARN=magenta ./myapp
```

### Programmatic Configuration

```go
log.SetSourceWidth(25)      // Set source column width
log.SetColorSource("cyan")  // Change source color
log.DisableColors()         // Disable all colors
```

## Output Format

Format: `[time] LEVEL [source] message key=value`

```
[2025/12/27 09:20:18] TRACE [main.go:16          ] trace message key=value
[2025/12/27 09:20:18] DEBUG [main.go:17          ] debug info user_id=123
[2025/12/27 09:20:18] INFO  [main.go:18          ] server started port=8080
[2025/12/27 09:20:18] WARN  [service.go:42       ] slow query duration=500ms
[2025/12/27 09:20:18] ERROR [exchange_client.go:1] connection failed err=timeout
```

- **Time**: in brackets `[YYYY/MM/DD HH:MM:SS]`
- **Level**: fixed 5-char width, colored
- **Source**: in brackets, green color by default, fixed width (default 20)
- **Message**: plain text with key=value pairs

## License

MIT
