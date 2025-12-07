// Package utils provides some utility functions that may be of use that don't fit
// any hard categories.
package utils

import (
	"fmt"
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
