package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/database"
	"luny.dev/cherryauctions/internal/routes"
)

func main() {
	db := database.SetupDatabase()
	database.MigrateModels(db)

	server := gin.New()
	routes.SetupRoutes(server, db)

	err := server.Run(":80")
	if err != nil {
		log.Fatalln("fatal: failed to run the server. conflicted port?")
	}
}
