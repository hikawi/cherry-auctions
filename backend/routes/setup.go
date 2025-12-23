package routes

import (
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	_ "luny.dev/cherryauctions/docs"
	"luny.dev/cherryauctions/routes/auth"
	"luny.dev/cherryauctions/routes/test"
	"luny.dev/cherryauctions/routes/users"
	"luny.dev/cherryauctions/utils"
)

type ServerDependency struct {
	Version    string
	DB         *gorm.DB
	S3Client   *s3.Client
	MailDialer *gomail.Dialer
}

func SetupServer(server *gin.Engine, db *gorm.DB) {
	err := server.SetTrustedProxies(nil)
	if err != nil {
		log.Println("warning: can't set trusted proxies", err.Error())
	}

	server.Use(gin.Recovery())
	server.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(utils.Getenv("CORS_ORIGINS", "http://localhost:5173"), ","),
		AllowMethods:     strings.Split(utils.Getenv("CORS_METHODS", "GET,HEAD,POST,PUT,DELETE"), ","),
		AllowCredentials: true,
		AllowHeaders:     strings.Split(utils.Getenv("CORS_HEADERS", ""), ","),
		AllowWebSockets:  true,
	}))
}

func SetupRoutes(server *gin.Engine, deps ServerDependency) {
	versionedGroup := server.Group(deps.Version)

	authHandler := auth.AuthHandler{DB: deps.DB}
	authHandler.SetupRouter(versionedGroup)

	usersHandler := users.UsersHandler{DB: deps.DB}
	usersHandler.SetupRouter(versionedGroup)

	testHandler := test.TestHandler{S3Client: deps.S3Client, MailDialer: deps.MailDialer}
	testHandler.SetupRouter(versionedGroup)

	versionedGroup.GET("/health", GetHealth)

	// Setup GIN swagger
	server.GET("/swagger", func(g *gin.Context) {
		g.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
