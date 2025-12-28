// Package closer provides a way to defer closing of resources, with checking for errors
// to appease the linter.
package closer

import (
	"fmt"
	"io"
)

func CloseResources(c io.Closer) {
	if err := c.Close(); err != nil {
		fmt.Println("warning: error while closing resource, error:", err)
	}
}
