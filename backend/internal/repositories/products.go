package repositories

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		Joins("INNER JOIN bids ON products.current_highest_bid_id = bids.id").
		Preload("Seller").
		Preload("Categories").
		Preload("CurrentHighestBid").
		Where("products.expired_at > ?", time.Now()).
		Order("bids.price DESC").
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
		Preload("DescriptionChanges", nil).
		Preload("DeniedBidders.User", nil).
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

// GetRunningUserProducts returns a list of products for a user, paged.
func (r *ProductRepository) GetRunningUserProducts(ctx context.Context, userID uint, limit int, offset int) ([]models.Product, error) {
	return gorm.G[models.Product](r.DB).
		Preload("Seller", nil).
		Preload("Categories", nil).
		Preload("CurrentHighestBid", nil).
		Where("seller_id = ? AND expired_at > ?", userID, time.Now()).
		Order("expired_at").
		Limit(limit).
		Offset(offset).
		Find(ctx)
}

func (r *ProductRepository) CountRunningUserProducts(ctx context.Context, userID uint) (int64, error) {
	return gorm.G[models.Product](r.DB).
		Where("seller_id = ? AND expired_at > ?", userID, time.Now()).
		Count(ctx, "id")
}

// CreateProduct creates a new product.
func (r *ProductRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	return gorm.G[models.Product](r.DB).Create(ctx, product)
}

// CreateDescriptionChange makes a new change to the description.
func (r *ProductRepository) CreateDescriptionChange(ctx context.Context, productID uint, productDescription string) error {
	return r.DB.WithContext(ctx).
		Model(&models.Product{Model: gorm.Model{ID: productID}}).
		Association("DescriptionChanges").
		Append(&models.DescriptionChange{Changes: productDescription})
}

// CreateBid creates a bid and returns the corresponding results.
func (r *ProductRepository) CreateBid(
	ctx context.Context,
	productID uint,
	userID uint,
	bidAmount int64,
	lastBid *models.Bid,
	newBid *models.Bid,
	currentProduct *models.Product,
) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Retrieve the product about to update
		product := models.Product{}
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&models.Product{}).
			Preload("CurrentHighestBid.User").
			Preload("Seller").
			Where("id = ?", productID).
			First(&product).
			Error
		if err != nil {
			return err
		}

		// Check if it is the highest price.
		// Whether it makes sense (based on step bid is on the service, not this I think)
		if (product.CurrentHighestBid != nil && product.CurrentHighestBid.Price >= bidAmount) || bidAmount < product.StartingBid {
			return fmt.Errorf("bid is not high enough")
		}

		// Check if it's the same user.
		if product.CurrentHighestBid != nil && product.CurrentHighestBid.UserID == userID {
			return fmt.Errorf("you can't outbid yourself")
		}

		// Now try to insert into the bids table I guess.
		bid := models.Bid{Price: bidAmount, UserID: userID, ProductID: product.ID}
		err = tx.Model(&bid).Create(&bid).Error
		if err != nil {
			return err
		}

		if product.CurrentHighestBid != nil {
			*lastBid = *product.CurrentHighestBid
		}
		*currentProduct = product
		*newBid = bid

		// Now mark the bids as counted.
		return tx.Model(&models.Product{Model: gorm.Model{ID: productID}}).
			Select("current_highest_bid_id", "bids_count").
			Updates(map[string]any{"current_highest_bid_id": bid.ID, "bids_count": gorm.Expr("bids_count + 1")}).
			Error
	})
}
