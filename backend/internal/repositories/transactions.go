package repositories

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"luny.dev/cherryauctions/internal/models"
)

type TransactionRepository struct {
	db          *gorm.DB
	productRepo *ProductRepository
	ratingRepo  *RatingRepostory
}

func NewTransactionRepository(
	db *gorm.DB,
	productRepo *ProductRepository,
	ratingRepo *RatingRepostory,
) *TransactionRepository {
	return &TransactionRepository{
		db:          db,
		productRepo: productRepo,
		ratingRepo:  ratingRepo,
	}
}

// CreateTransaction creates a new transaction object bound to a certain product.
func (r *TransactionRepository) CreateTransaction(ctx context.Context, transaction *models.Transaction) error {
	return r.db.WithContext(ctx).
		Model(&models.Transaction{}).
		Create(transaction).
		Error
}

func (r *TransactionRepository) GetTransactionByID(ctx context.Context, id uint) (models.Transaction, error) {
	trans := models.Transaction{}
	err := r.db.WithContext(ctx).
		Model(&models.Transaction{}).
		Preload("Product").
		Preload("Seller").
		Preload("Buyer").
		Where("id = ?", id).
		First(&trans).
		Error
	return trans, err
}

func (r *TransactionRepository) UpdateTransactionStatus(ctx context.Context, id uint, status models.TransactionStatus) (int64, error) {
	db := r.db.WithContext(ctx).
		Model(&models.Transaction{}).
		Where("id = ?", id).
		Update("transaction_status", status)
	return db.RowsAffected, db.Error
}

func (r *TransactionRepository) CancelTransactionStatus(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var transaction models.Transaction
		err := tx.Model(&models.Transaction{}).
			Where("id = ?", id).
			First(&transaction).
			Error
		if err != nil {
			return err
		}

		// Doesn't allow the thing to be cancelled if the winner has already paid.
		if transaction.TransactionStatus != models.TransactionStatusPending {
			return fmt.Errorf("transaction %d is not pending", id)
		}

		db := tx.Model(&models.Transaction{}).
			Where("id = ?", id).
			Update("transaction_status", models.TransactionStatusCancelled)
		if db.Error != nil || db.RowsAffected == 0 {
			return fmt.Errorf("transaction could not be set to cancelled")
		}

		// Mark the winner as bad.
		rating := models.Rating{
			ProductID:  transaction.ProductID,
			ReviewerID: transaction.SellerID,
			RevieweeID: transaction.SellerID,
			Rating:     0,
			Feedback:   "Did not follow through with payment",
		}
		err = r.ratingRepo.CreateRating(ctx, &rating)
		if err != nil {
			return err
		}

		return nil
	})
}
