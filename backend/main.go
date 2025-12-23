package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/database"
	"luny.dev/cherryauctions/routes"
	"luny.dev/cherryauctions/services"
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

	s3Client := services.NewS3Service()
	log.Println(*s3Client.Options().BaseEndpoint)
	output, err := s3Client.ListBuckets(context.Background(), &s3.ListBucketsInput{
		BucketRegion: aws.String("us-east-1"),
	})
	if err != nil {
		log.Fatal(err)
	}

	for bucket := range output.Buckets {
		log.Println(bucket)
	}

	server := gin.New()

	routes.SetupServer(server, db)
	routes.SetupRoutes(server, db)

	err = server.Run(":80")
	if err != nil {
		log.Fatalln("fatal: failed to run the server. conflicted port?")
	}
}
