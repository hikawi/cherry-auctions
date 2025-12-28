package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type LogLevel string

const (
	LOG_DEBUG LogLevel = "debug"
	LOG_INFO  LogLevel = "info"
	LOG_WARN  LogLevel = "warn"
	LOG_ERROR LogLevel = "error"
)

var logger io.Writer

func InitLogger() {
	if logger != nil {
		return
	}

	f, err := os.OpenFile("./server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("fatal: can't open file to log %s", err)
	}

	logger = io.MultiWriter(f)
}

// LogMessage logs a message with a certain log level.
func LogMessage(c *gin.Context, level LogLevel, message map[string]any) {
	// Remove password-related headers.
	headers := c.Request.Header.Clone()
	headers.Del("Authorization")
	headers.Del("Cookie")

	jsonMsg := map[string]any{
		"ip":        c.ClientIP(),
		"path":      c.Request.URL.Path,
		"method":    c.Request.Method,
		"headers":   headers,
		"log_level": level,
		"timestamp": time.Now(),
	}

	// Combine keys.
	maps.Copy(jsonMsg, message)

	// Marshal and write
	val, err := json.Marshal(jsonMsg)
	if err != nil {
		fmt.Printf("warning: unable to marshal log %s\n", message)
	}

	bytes, err := fmt.Fprintf(logger, "%s\n", val)
	if err != nil || bytes == 0 {
		fmt.Printf("warning: unable to write log %s", err)
	}
}
