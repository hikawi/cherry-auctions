package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/routes"
)

const version string = "v1"

func main() {
	server := gin.Default()

	versionedGroup := server.Group(version)
	versionedGroup.GET("/health", routes.GetHealth)

	err := server.Run(":80")
	if err != nil {
		log.Fatalln("fatal: failed to run the server. conflicted port?")
	}
}
