package main

import (
	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/routes"
)

const version string = "v1"

func main() {
	server := gin.Default()

	versionedGroup := server.Group(version)
	versionedGroup.GET("/health", routes.GetHealth)

	server.Run(":80")
}
