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
    log.Init()          // LOG_LEVEL from env (default: INFO)
    
    // Optional: set source location width (default: 15)
    log.SetSourceWidth(12)
    
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

Set `LOG_LEVEL` environment variable:

```bash
# Development
LOG_LEVEL=DEBUG ./myapp

# Production
LOG_LEVEL=INFO ./myapp
```

## Output Example

Format: `[time] LEVEL [source] message key=value`

```
[2025/12/26 16:00:00] TRACE [main.go:16     ] trace message key=value
[2025/12/26 16:00:00] DEBUG [main.go:17     ] debug info user_id=123
[2025/12/26 16:00:00] INFO  [main.go:18     ] server started port=8080
[2025/12/26 16:00:00] WARN  [main.go:19     ] slow query duration=500ms
[2025/12/26 16:00:00] ERROR [main.go:20     ] connection failed err=timeout
```

- **Time**: in brackets
- **Level**: fixed 5-char width, colored (TRACE=light gray, DEBUG=gray, WARN=yellow, ERROR/FATAL/PANIC=red)
- **Source**: in brackets, green color, fixed width (default 15, configurable via `SetSourceWidth`)
- **Message**: plain text with key=value pairs

## License

MIT
