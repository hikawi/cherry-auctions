package repositories

import (
	"context"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db          *gorm.DB
	productRepo *ProductRepository
}

func NewTransactionRepository(
	db *gorm.DB,
	productRepo *ProductRepository,
) *TransactionRepository {
	return &TransactionRepository{
		db:          db,
		productRepo: productRepo,
	}
}

// CreateTransaction creates a new transaction object bound to a certain product.
func (r *TransactionRepository) CreateTransaction(ctx context.Context, productID uint) {
}
