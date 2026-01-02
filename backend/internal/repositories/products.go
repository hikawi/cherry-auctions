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

func (r *ProductRepository) GetProductByID(ctx context.Context, id int) (models.Product, error) {
	return gorm.G[models.Product](r.DB).
		Preload("Seller", nil).
		Preload("Categories", nil).
		Preload("CurrentHighestBid.User", nil).
		Preload("Bids.User", nil).
		Preload("Questions.User", nil).
		Preload("Categories", nil).
		Preload("ProductImages", nil).
		Where("id = ?", id).
		First(ctx)
}

func (r *ProductRepository) GetSimilarProductsTo(ctx context.Context, product *models.Product) ([]models.Product, error) {
	var categoryIDs []uint
	for _, cat := range product.Categories {
		categoryIDs = append(categoryIDs, cat.ID)
	}

	var products []models.Product
	err := r.DB.Model(&models.Product{}).WithContext(ctx).
		Joins("JOIN products_categories ON products_categories.product_id = products.id").
		Where("products.id <> ?", product.ID).                      // Exclude current product
		Where("products_categories.category_id IN ?", categoryIDs). // Match any shared category
		Preload("Seller").
		Preload("Categories").
		Preload("CurrentHighestBid").
		Distinct("products.*").
		Limit(5).
		Find(&products).
		Error
	return products, err
}

func (r *ProductRepository) GetFavoriteProducts(ctx context.Context, userID uint, limit int, offset int) ([]models.Product, error) {
	var products []models.Product
	err := r.DB.Model(&models.Product{}).WithContext(ctx).
		Joins("JOIN favorite_products ON products.id = favorite_products.product_id").
		Where("favorite_products.user_id = ?", userID).
		Preload("Seller").
		Preload("Categories").
		Preload("CurrentHighestBid").
		Distinct("products.*").
		Limit(limit).
		Offset(offset).
		Order("products.id").
		Find(&products).
		Error
	return products, err
}

func (r *ProductRepository) CountFavoriteProducts(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.DB.Table("favorite_products").Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *ProductRepository) ToggleFavoriteProduct(ctx context.Context, userID uint, productID uint) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		data := models.FavoriteProduct{
			UserID:    userID,
			ProductID: productID,
		}
		_, err := gorm.G[models.FavoriteProduct](tx).Where("user_id = ? AND product_id = ?", userID, productID).First(ctx)
		if err == nil {
			_, err = gorm.G[models.FavoriteProduct](tx).Where("user_id = ? AND product_id = ?", userID, productID).Delete(ctx)
			return err
		}

		return gorm.G[models.FavoriteProduct](tx).Create(ctx, &data)
	})
}

// AttachFavoriteStatus is a temporary function to attach a "IsFavorite" field
// to products, just so we don't have to make too many queries.
func (r *ProductRepository) AttachFavoriteStatus(ctx context.Context, userID uint, products ...*models.Product) {
	if userID == 0 || len(products) == 0 {
		return
	}

	// 1. Collect all IDs from the products we just fetched
	ids := make([]uint, len(products))
	for i, p := range products {
		ids[i] = p.ID
	}

	// 2. Query the favorite table once for the whole batch
	var favoriteIDs []uint
	r.DB.WithContext(ctx).Table("favorite_products").
		Where("user_id = ? AND product_id IN ?", userID, ids).
		Pluck("product_id", &favoriteIDs)

	// 3. Map the results back to the structs
	favMap := make(map[uint]bool)
	for _, id := range favoriteIDs {
		favMap[id] = true
	}

	for _, p := range products {
		p.IsFavorite = favMap[p.ID]
	}
}

// CreateProduct creates a new product.
func (r *ProductRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	return gorm.G[models.Product](r.DB).Create(ctx, product)
}
