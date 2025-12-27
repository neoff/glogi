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

```
2025/12/26 16:00:00 main.go:16: [TRACE] trace message key=value
2025/12/26 16:00:00 main.go:17: [DEBUG] debug info user_id=123
2025/12/26 16:00:00 main.go:18: [INFO] server started port=8080
2025/12/26 16:00:00 main.go:19: [WARN] slow query duration=500ms
2025/12/26 16:00:00 main.go:20: [ERROR] connection failed err=timeout
```

## License

MIT
