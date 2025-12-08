// Package utils provides some utility functions that may be of use that don't fit
// any hard categories.
package utils

import (
	"fmt"
	"log"
	"os"
)

// Getenv retrieves an environment value and uses a default value if not available.
func Getenv(key string, def string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		fmt.Printf("warning: unable to find environment variable for key = %s\n", key)
		return def
	}

	return val
}

// Fatalenv retrieves an environment value and kills itself if it doesn't exist.
func Fatalenv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("fatal: unable to find required environment variable for key %s\n", key)
	}

	return val
}
