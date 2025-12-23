package utils

import (
	"io"

	"github.com/gin-gonic/gin"
)

func CloseResources(c io.Closer) {
	if err := c.Close(); err != nil {
		Log(gin.H{"message": "error while closing resource", "error": err})
	}
}
