package repositories

import (
	"context"

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

func (r *ProductRepository) CountProducts(ctx context.Context) (int64, error) {
	return gorm.G[models.Product](r.DB).Count(ctx, "id")
}
