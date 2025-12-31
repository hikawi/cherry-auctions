package repositories

import (
	"context"
	"time"

	"gorm.io/gorm"
	"luny.dev/cherryauctions/internal/models"
)

type ProductRepository struct {
	DB *gorm.DB
}

func (r *ProductRepository) SearchProducts(ctx context.Context, query string, limit int, offset int) ([]models.Product, error) {
	statement := gorm.G[models.Product](r.DB).Preload("Seller", nil)

	// Conditionally apply full-text and fuzzy.
	if query != "" {
		return statement.Where("(search_vector @@ plainto_tsquery('simple', ?)) OR name % ?", query, query).
			Order(gorm.Expr(
				"(ts_rank(search_vector, plainto_tsquery('simple', ?)) * 2.0) + similarity(name, ?) DESC", // Just weigh the full-text better
				query, query,
			)).
			Limit(limit).
			Offset(offset).
			Find(ctx)
	}

	return statement.Limit(limit).Offset(offset).Find(ctx)
}

func (r *ProductRepository) CountProductsWithQuery(ctx context.Context, query string) (int64, error) {
	statement := gorm.G[models.Product](r.DB).Preload("Seller", nil)

	// Conditionally apply full-text and fuzzy.
	if query != "" {
		return statement.Where("(search_vector @@ plainto_tsquery('simple', ?)) OR name % ?", query, query).
			Order(gorm.Expr(
				"(ts_rank(search_vector, plainto_tsquery('simple', ?)) * 2.0) + similarity(name, ?) DESC", // Just weigh the full-text better
				query, query,
			)).
			Count(ctx, "id")
	}

	return statement.Count(ctx, "id")
}

func (r *ProductRepository) CountProducts(ctx context.Context) (int64, error) {
	return gorm.G[models.Product](r.DB).Count(ctx, "id")
}

// GetTopEndingSoons returns 5 products that are currently about to expire.
func (r *ProductRepository) GetTopEndingSoons(ctx context.Context) ([]models.Product, error) {
	return gorm.G[models.Product](r.DB).
		Preload("Seller", nil).
		Preload("Categories", nil).
		Preload("CurrentHighestBid", nil).
		Where("expired_at > ?", time.Now()).
		Order("expired_at ASC").
		Limit(5).
		Find(ctx)
}

func (r *ProductRepository) GetMostActiveProducts(ctx context.Context) ([]models.Product, error) {
	return gorm.G[models.Product](r.DB).
		Preload("Seller", nil).
		Preload("Categories", nil).
		Preload("CurrentHighestBid", nil).
		Where("expired_at > ?", time.Now()).
		Order("bids_count DESC, expired_at ASC").
		Limit(5).
		Find(ctx)
}

func (r *ProductRepository) GetHighestBiddedProducts(ctx context.Context) ([]models.Product, error) {
	var products []models.Product

	err := r.DB.WithContext(ctx).
		Joins("LEFT JOIN bids ON products.current_highest_bid_id = bids.id").
		Preload("Seller").
		Preload("Categories").
		Preload("CurrentHighestBid").
		Where("expired_at > ?", time.Now()).
		Order("bids.price DESC, expired_at ASC").
		Limit(5).
		Find(&products).
		Error

	return products, err
}
