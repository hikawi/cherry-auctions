package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const version string = "v1"

func SetupRoutes(server *gin.Engine, db *gorm.DB) {
	server.Use(gin.Recovery())
	server.SetTrustedProxies(nil)

	versionedGroup := server.Group(version)
	versionedGroup.GET("/health", GetHealth)
}
