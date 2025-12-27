package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/routes"
	"luny.dev/cherryauctions/services"
	"luny.dev/cherryauctions/utils"
)

// @title						Cherry Auctions API
// @version					1.0
// @description				Backend API for CherryAuctions at cherry-auctions.luny.dev.
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
// @description				Classic Bearer token, authenticated by using the login endpoint, which should grant an access token. To refresh it, use the RefreshToken cookie.
func main() {
	utils.InitLogger()

	db := services.SetupDatabase()
	s3Client := services.NewS3Service()
	mailDialer := services.NewMailerService()

	// Weird to do this even in production.
	services.MigrateModels(db)

	server := gin.New()

	routes.SetupServer(server, db)
	routes.SetupRoutes(server, routes.ServerDependency{
		Version:    "v1",
		DB:         db,
		S3Client:   s3Client,
		MailDialer: mailDialer,
	})

	err := server.Run(":80")
	if err != nil {
		log.Fatalln("fatal: failed to run the server. conflicted port?")
	}
}
