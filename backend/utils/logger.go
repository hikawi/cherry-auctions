package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

// Log logs a piece of JSON data. Stub for when OpenSearch is fully implemented.
func Log(data any) {
	val, err := json.Marshal(data)
	time := time.Now()

	if err != nil {
		fmt.Printf("%s | unable to marshal log %s\n", time.Format("hh:mm:ss"), data)
	}

	fmt.Printf("%s | %s", time.Format("hh:mm:ss"), val)
}
