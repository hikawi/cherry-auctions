package routes

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	_ "luny.dev/cherryauctions/docs"
	"luny.dev/cherryauctions/routes/auth"
	"luny.dev/cherryauctions/utils"
)

const version string = "v1"

func SetupServer(server *gin.Engine, db *gorm.DB) {
	server.SetTrustedProxies(nil)
	server.Use(gin.Recovery())
	server.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(utils.Getenv("CORS_ORIGINS", "http://localhost:5173"), ","),
		AllowMethods:     strings.Split(utils.Getenv("CORS_METHODS", "GET,HEAD,POST,PUT,DELETE"), ","),
		AllowCredentials: true,
		AllowHeaders:     strings.Split(utils.Getenv("CORS_HEADERS", ""), ","),
		AllowWebSockets:  true,
	}))
}

func SetupRoutes(server *gin.Engine, db *gorm.DB) {
	versionedGroup := server.Group(version)

	authHandler := auth.AuthHandler{DB: db}
	authHandler.SetupRouter(versionedGroup)

	versionedGroup.GET("/health", GetHealth)

	// Setup GIN swagger
	server.GET("/swagger", func(g *gin.Context) {
		g.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
