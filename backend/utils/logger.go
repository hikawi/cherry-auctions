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

	// Go uses weird ass formatting for time... Refer to their docs.
	// 15 == hh
	// 04 == mm
	// 05 == ss
	if err != nil {
		fmt.Printf("%s | unable to marshal log %s\n", time.Format("15:04:05"), data)
	}
	fmt.Printf("%s | %s", time.Format("15:04:05"), val)
}
