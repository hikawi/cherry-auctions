package repositories

import (
	"context"

	"gorm.io/gorm"
	"luny.dev/cherryauctions/internal/models"
)

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

func (r *ChatSessionRepository) GetUserChatSessions(ctx context.Context, userID uint, limit int, offset int) ([]models.ChatSession, error) {
	var sessions []models.ChatSession
	err := r.db.WithContext(ctx).
		Model(&models.ChatSession{}).
		Preload("Seller").
		Preload("Buyer").
		Preload("Product").
		Where("buyer_id = ? OR seller_id = ?", userID, userID).
		Order("updated_at DESC, created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&sessions).
		Error
	return sessions, err
}

func (r *ChatSessionRepository) CountUserChatSessions(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.ChatSession{}).
		Where("buyer_id = ? OR seller_id = ?", userID, userID).
		Count(&count).
		Error
	return count, err
}

func (r *ChatSessionRepository) CreateChatSession(ctx context.Context, session *models.ChatSession) error {
	return r.db.WithContext(ctx).
		Model(&models.ChatSession{}).
		Create(session).
		Error
}

func (r *ChatSessionRepository) GetChatSessionByID(ctx context.Context, id uint) (models.ChatSession, error) {
	session := models.ChatSession{}
	err := r.db.WithContext(ctx).
		Model(&models.ChatSession{}).
		Where("id = ?", id).
		First(&session).
		Error
	return session, err
}

func (r *ChatSessionRepository) GetSessionChatMessages(ctx context.Context, sessionID uint, limit int, offset int) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := r.db.WithContext(ctx).
		Model(&models.ChatMessage{}).
		Preload("Sender").
		Where("chat_session_id = ?", sessionID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).
		Error
	return messages, err
}

func (r *ChatSessionRepository) CountSessionChatMessages(ctx context.Context, sessionID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.ChatMessage{}).
		Where("chat_session_id = ?", sessionID).
		Count(&count).
		Error
	return count, err
}

func (r *ChatSessionRepository) CreateChatMessage(ctx context.Context, msg *models.ChatMessage) error {
	return r.db.WithContext(ctx).
		Model(&models.ChatMessage{}).
		Create(msg).
		Error
}
