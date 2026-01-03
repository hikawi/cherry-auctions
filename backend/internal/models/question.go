package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type Question struct {
	gorm.Model
	Content string `gorm:"not null"`
	Answer  sql.NullString

	ProductID uint
	Product   Product
	UserID    uint
	User      User // Who asked
}
