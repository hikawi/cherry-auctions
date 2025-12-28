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
	"luny.dev/cherryauctions/internal/config"
	"luny.dev/cherryauctions/internal/repositories"
	"luny.dev/cherryauctions/internal/routes/auth"
	"luny.dev/cherryauctions/internal/routes/categories"
	"luny.dev/cherryauctions/internal/routes/users"
	"luny.dev/cherryauctions/internal/services"
	"luny.dev/cherryauctions/pkg/env"
)

type ServerDependency struct {
	Version      string
	DB           *gorm.DB
	S3Client     *s3.Client
	MailDialer   *gomail.Dialer
	Config       *config.Config
	Services     services.ServiceRegistry
	Repositories repositories.RepositoryRegistry
}

func SetupServer(server *gin.Engine, db *gorm.DB) {
	err := server.SetTrustedProxies(nil)
	if err != nil {
		log.Println("warning: can't set trusted proxies", err.Error())
	}

	server.Use(gin.Recovery())
	server.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(env.Getenv("CORS_ORIGINS", "http://localhost:5173"), ","),
		AllowMethods:     strings.Split(env.Getenv("CORS_METHODS", "GET,HEAD,POST,PUT,DELETE"), ","),
		AllowCredentials: true,
		AllowHeaders:     strings.Split(env.Getenv("CORS_HEADERS", "Authorization"), ","),
		AllowWebSockets:  true,
	}))
}

func SetupRoutes(server *gin.Engine, deps ServerDependency) {
	versionedGroup := server.Group(deps.Version)

	authHandler := auth.AuthHandler{
		DB:               deps.DB,
		CookieSecure:     deps.Config.CookieSecure,
		Domain:           deps.Config.Domain,
		JWTService:       deps.Services.JWTService,
		RandomService:    deps.Services.RandomService,
		PasswordService:  deps.Services.PasswordService,
		CaptchaService:   deps.Services.CaptchaService,
		UserRepo:         deps.Repositories.UserRepository,
		RefreshTokenRepo: deps.Repositories.RefreshTokenRepository,
	}
	authHandler.SetupRouter(versionedGroup)

	usersHandler := users.UsersHandler{
		DB:                deps.DB,
		MiddlewareService: deps.Services.MiddlewareService,
		UserRepo:          deps.Repositories.UserRepository,
	}
	usersHandler.SetupRouter(versionedGroup)

	categoriesHandler := categories.CategoriesHandler{
		DB:                deps.DB,
		MiddlewareService: deps.Services.MiddlewareService,
		CategoryRepo:      deps.Repositories.CategoryRepository,
	}
	categoriesHandler.SetupRouter(versionedGroup)

	versionedGroup.GET("/health", GetHealth)

	// Setup GIN swagger
	server.GET("/swagger", func(g *gin.Context) {
		g.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
