package utils

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

func LogMessage(c *gin.Context, level LogLevel, message map[string]any) {
	// Remove password-related headers.
	headers := c.Request.Header
	headers.Del("Authorization")

	jsonMsg := map[string]any{
		"ip":        c.ClientIP(),
		"path":      c.Request.URL.Path,
		"method":    c.Request.Method,
		"headers":   headers,
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

// Log logs a piece of JSON data. Stub for when OpenSearch is fully implemented.
// Deprecated!
func Log(data any) {
	val, err := json.Marshal(data)
	time := time.Now()

	// Go uses weird ass formatting for time... Refer to their docs.
	// 15 == hh
	// 04 == mm
	// 05 == ss
	if err != nil {
		fmt.Printf("%s | unable to marshal log %s\n", time.Format("15:04:05"), data)
	} else {
		fmt.Printf("%s | %s\n", time.Format("15:04:05"), val)
	}
}
