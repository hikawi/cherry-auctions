package repositories

import "gorm.io/gorm"

type ChatSessionRepository struct {
	db *gorm.DB
}

func NewChatSessionRepository(
	db *gorm.DB,
) *ChatSessionRepository {
	return &ChatSessionRepository{
		db: db,
	}
}
