package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/database"
	"luny.dev/cherryauctions/routes"
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
// @schemes					http https
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Classic Bearer token
func main() {
	db := database.SetupDatabase()
	database.MigrateModels(db)

	server := gin.New()

	routes.SetupServer(server, db)
	routes.SetupRoutes(server, db)

	err := server.Run(":80")
	if err != nil {
		log.Fatalln("fatal: failed to run the server. conflicted port?")
	}
}
