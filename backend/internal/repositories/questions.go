package repositories

import (
	"context"

	"gorm.io/gorm"
	"luny.dev/cherryauctions/internal/models"
)

type QuestionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{
		db: db,
	}
}

func (r *QuestionRepository) GetQuestionByID(ctx context.Context, id uint) (models.Question, error) {
	return gorm.G[models.Question](r.db).
		Preload("User", nil).
		Preload("Product.Bids.User", nil).
		Where("id = ?", id).
		First(ctx)
}

// CreateQuestion creates a question under a product
func (r *QuestionRepository) CreateQuestion(ctx context.Context, productID uint, userID uint, content string) error {
	return r.db.WithContext(ctx).
		Model(&models.Product{Model: gorm.Model{ID: productID}}).
		Association("Questions").
		Append(&models.Question{UserID: userID, Content: content})
}

// AnswerProductQuestion appends an answer to a product.
func (r *QuestionRepository) AnswerProductQuestion(ctx context.Context, questionID uint, answer string) (int, error) {
	return gorm.G[models.Question](r.db).
		Where("id = ?", questionID).
		Update(ctx, "answer", answer)
}
