package models

import "time"

const (
	ROLE_USER  = "user"
	ROLE_ADMIN = "admin"
)

type Role struct {
	ID          string `gorm:"primaryKey"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
