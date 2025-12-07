package main

import (
	"log"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/database"
	"luny.dev/cherryauctions/internal/routes"
	"luny.dev/cherryauctions/internal/utils"
)

// @title						Cherry Auctions API
// @version					0.1
// @description				Backend API for CherryAuctions.
// @contact.name				Nguyệt Ánh
// @contact.email				hello@luny.dev
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @basepath					/v1
// @accept						json
// @produce					json
// @schemes					https
// @securityDefinitions.apikey	ApiKeyAuth
// @in header
// @name Authorization
// @description Classic Bearer token
func main() {
	db := database.SetupDatabase()
	database.MigrateModels(db)

	server := gin.New()
	server.Use(gin.Recovery())
	server.Use(gin.Logger())
	server.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(utils.Getenv("CORS_ORIGINS", "http://localhost:5173"), ","),
		AllowMethods:     strings.Split(utils.Getenv("CORS_METHODS", "GET,HEAD,POST,PUT,DELETE"), ","),
		AllowCredentials: true,
		AllowHeaders:     strings.Split(utils.Getenv("CORS_HEADERS", ""), ","),
		AllowWebSockets:  true,
	}))

	routes.SetupRoutes(server, db)

	err := server.Run(":80")
	if err != nil {
		log.Fatalln("fatal: failed to run the server. conflicted port?")
	}
}
