package main

import (
	"errors"

	log "github.com/neoff/glogi"
)

func main() {
	defer log.Recover()

	log.Init()

	// Structured logging
	log.Trace("trace message", "key", "value")
	log.Debug("debug info", "user_id", 123)
	log.Info("server started", "port", 8080)
	log.Warn("slow query", "duration", "500ms")

	// Error handling
	err := errors.New("connection timeout")
	log.Error("database error", "err", err)

	// Standard log compatibility
	log.Println("This is a println message")
	log.Printf("Formatted message: %d", 42)

	log.Info("all tests passed")
}
