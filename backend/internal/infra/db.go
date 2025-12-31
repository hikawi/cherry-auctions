// Package infra provides a bunch of functions to setup basic raw connections.
// This should be piped directly to handlers.
package infra

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"luny.dev/cherryauctions/internal/models"
)

func SetupDatabase(databaseURL string) *gorm.DB {
	conn, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
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
	err := db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.Category{},
		&models.Role{},
		&models.Product{},
		&models.Question{},
		&models.ProductImage{},
		&models.SellerSubscription{},
		&models.Bid{},
	)
	if err != nil {
		log.Fatalln("fatal: failed to auto migrate models. check them yourself")
	}
}
