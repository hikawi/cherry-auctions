package services

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"luny.dev/cherryauctions/models"
)

// SetupDatabase sets up the database, and returns the *gorm.DB
func SetupDatabase() *gorm.DB {
	conn, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatalln("fatal: can't connect to database")
	}

	db, err := conn.DB()
	if err != nil {
		log.Fatalln("fatal: can't setup database")
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(30)
	db.SetConnMaxIdleTime(time.Minute * 15)
	db.SetConnMaxLifetime(time.Hour)

	return conn
}

// MigrateModels uses GORM to migrate the models.
func MigrateModels(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{}, &models.RefreshToken{}, &models.Category{})
	if err != nil {
		log.Fatalln("fatal: failed to auto migrate models. check them yourself")
	}
}
