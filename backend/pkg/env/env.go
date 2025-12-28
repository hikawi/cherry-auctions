// Package env provides some utility functions that may be of use that don't fit
// any hard categories.
package env

import (
	"fmt"
	"log"
	"os"
	"strconv"
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

// FatalenvInt retrieves an environment value, kills itself if it doesn't exist OR if it is not a valid number.
func FatalenvInt(key string) int64 {
	val := Fatalenv(key)
	parsed, err := strconv.ParseInt(val, 10, 0)
	if err != nil {
		log.Fatalf("fatal: found env variable for %s but is not an integer: %s\n", key, val)
	}

	return parsed
}

// FatalenvBool retrieves an environment value, kills itself if it doesn't exist OR if it is not a valid boolean.
func FatalenvBool(key string) bool {
	val := Fatalenv(key)
	parsed, err := strconv.ParseBool(val)
	if err != nil {
		log.Fatalf("fatal: found env variable for %s but is not a boolean: %s\n", key, val)
	}

	return parsed
}
