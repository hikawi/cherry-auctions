package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/config"
	"luny.dev/cherryauctions/internal/infra"
	"luny.dev/cherryauctions/internal/logging"
	"luny.dev/cherryauctions/internal/repositories"
	"luny.dev/cherryauctions/internal/routes"
	"luny.dev/cherryauctions/internal/services"
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
	cfg := config.Load()

	logging.InitLogger()

	db := infra.SetupDatabase(cfg.DatabaseURL)
	s3Client := infra.SetupS3(cfg.AWS.S3Base, cfg.AWS.S3UsePathStyle)
	mailDialer := infra.SetupMailer(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.User, cfg.SMTP.Password)

	// Setup repositories here
	categoryRepo := &repositories.CategoryRepository{DB: db}
	roleRepo := &repositories.RoleRepository{DB: db}
	userRepo := &repositories.UserRepository{DB: db, RoleRepository: roleRepo}
	refreshTokenRepo := &repositories.RefreshTokenRepository{DB: db}
	productRepo := &repositories.ProductRepository{DB: db}

	// Setup services here
	jwtService := &services.JWTService{JWTDomain: cfg.Domain, JWTAudience: cfg.JWT.Audience, JWTSecretKey: cfg.JWT.Secret, JWTExpiry: cfg.JWT.Expiry}
	randomService := &services.RandomService{}
	passwordService := &services.PasswordService{RandomService: randomService}
	captchaService := &services.CaptchaService{RecaptchaSecret: cfg.RecaptchaSecret}
	middlewareService := &services.MiddlewareService{JWTService: jwtService}

	// Weird to do this even in production.
	infra.MigrateModels(db)

	server := gin.New()

	routes.SetupServer(server, db)
	routes.SetupRoutes(server, routes.ServerDependency{
		Version:    "v1",
		DB:         db,
		S3Client:   s3Client,
		MailDialer: mailDialer,
		Config:     cfg,
		Services: services.ServiceRegistry{
			JWTService:        jwtService,
			RandomService:     randomService,
			PasswordService:   passwordService,
			CaptchaService:    captchaService,
			MiddlewareService: middlewareService,
		},
		Repositories: repositories.RepositoryRegistry{
			CategoryRepository:     categoryRepo,
			UserRepository:         userRepo,
			RoleRepository:         roleRepo,
			RefreshTokenRepository: refreshTokenRepo,
			ProductRepository:      productRepo,
		},
	})

	err := server.Run(":80")
	if err != nil {
		log.Fatalln("fatal: failed to run the server. conflicted port?")
	}
}
