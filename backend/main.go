package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
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
	// Setup a cron tab
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("failed to setup scheduler: %v", err)
	}

	cfg := config.Load()
	vips.Startup(nil)
	defer vips.Shutdown()

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
	questionRepo := repositories.NewQuestionRepository(db)
	chatSessionRepo := repositories.NewChatSessionRepository(db)
	ratingRepo := repositories.NewRatingRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db, productRepo, ratingRepo)

	// Setup services here
	jwtService := &services.JWTService{JWTDomain: cfg.Domain, JWTAudience: cfg.JWT.Audience, JWTSecretKey: cfg.JWT.Secret, JWTExpiry: cfg.JWT.Expiry}
	randomService := &services.RandomService{}
	passwordService := &services.PasswordService{RandomService: randomService}
	captchaService := &services.CaptchaService{RecaptchaSecret: cfg.RecaptchaSecret}
	middlewareService := &services.MiddlewareService{JWTService: jwtService}
	s3Service := services.NewS3Service(cfg.AWS.BucketName, s3Client)
	mailerService := services.NewMailerService(cfg, mailDialer, productRepo, questionRepo, userRepo)
	otpService := services.NewOTPService(mailerService, userRepo)

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
			S3Service:         s3Service,
			MailerService:     mailerService,
			OTPService:        otpService,
		},
		Repositories: repositories.RepositoryRegistry{
			CategoryRepository:     categoryRepo,
			UserRepository:         userRepo,
			RoleRepository:         roleRepo,
			RefreshTokenRepository: refreshTokenRepo,
			ProductRepository:      productRepo,
			QuestionRepository:     questionRepo,
			ChatSessionRepository:  chatSessionRepo,
			TransactionRepository:  transactionRepo,
			RatingRepostory:        ratingRepo,
		},
	})

	_, err = scheduler.NewJob(gocron.DurationJob(1*time.Minute), gocron.NewTask(func() {
		// Sweep
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		err := productRepo.UpdateAllExpiredProducts(ctx)
		if err != nil {
			fmt.Printf("warning: unable to update expired products: %v\n", err)
		}

		mailerService.SendEndedAuctionsEmail()
	}))
	if err != nil {
		log.Fatalf("can't setup a cron job: %v", err)
	}
	scheduler.Start()
	defer func() {
		err := scheduler.Shutdown()
		if err != nil {
			log.Printf("couldn't shutdown scheduler: %v", err)
		}
	}()

	err = server.Run(":80")
	if err != nil {
		log.Fatalln("fatal: failed to run the server. conflicted port?")
	}
}
