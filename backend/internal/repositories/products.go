package repositories

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"luny.dev/cherryauctions/internal/models"
)

type ProductSortType string

const (
	ProductSortTypeNone       ProductSortType = "id"
	ProductSortTypeExpiryTime ProductSortType = "time"
	ProductSortTypePrice      ProductSortType = "price"
)

type ProductRepository struct {
	DB *gorm.DB
}

func (r *ProductRepository) SearchProducts(
	ctx context.Context,
	query string,
	categories []uint,
	sortType ProductSortType,
	sortAsc bool,
	limit int,
	offset int,
) ([]models.Product, error) {
	db := r.DB.WithContext(ctx).
		Model(&models.Product{}).
		Preload("Seller").
		Preload("Categories").
		Preload("CurrentHighestBid.User")

	if len(categories) > 0 {
		db = db.
			Joins("JOIN products_categories ON products_categories.product_id = products.id").
			Where("products_categories.category_id IN ?", categories).
			Distinct()
	}

	if query != "" {
		db = db.Where(
			r.DB.Where("search_vector @@ plainto_tsquery('simple', ?)", query).
				Or("name % ?", query),
		)
		rankExpr := "(ts_rank(search_vector, plainto_tsquery('simple', ?)) * 2.0) + similarity(name, ?) DESC"
		db = db.Order(gorm.Expr(rankExpr, query, query))
	}

	// Lol
	switch sortType {
	case ProductSortTypeNone:
		if sortAsc {
			db = db.Order("products.id ASC")
		} else {
			db = db.Order("products.id DESC")
		}
	case ProductSortTypeExpiryTime:
		if sortAsc {
			db = db.Order("products.expired_at ASC")
		} else {
			db = db.Order("products.expired_at DESC")
		}
	case ProductSortTypePrice:
		if sortAsc {
			db = db.Order("products.bin_price ASC")
		} else {
			db = db.Order("products.bin_price DESC")
		}
	}

	var products []models.Product
	err := db.
		Where("product_state = ?", models.ProductStateActive).
		Limit(limit).
		Offset(offset).
		Find(&products).
		Error
	return products, err
}

func (r *ProductRepository) CountProductsWithQuery(ctx context.Context, query string, categories []uint) (int64, error) {
	db := r.DB.WithContext(ctx).
		Model(&models.Product{}).
		Preload("Seller").
		Preload("Categories")

	if len(categories) > 0 {
		db = db.
			Joins("JOIN products_categories ON products_categories.product_id = products.id").
			Where("products_categories.category_id IN ?", categories).
			Distinct()
	}

	if query != "" {
		db = db.Where(
			r.DB.Where("search_vector @@ plainto_tsquery('simple', ?)", query).
				Or("name % ?", query),
		)
		rankExpr := "(ts_rank(search_vector, plainto_tsquery('simple', ?)) * 2.0) + similarity(name, ?) DESC"
		db = db.Order(gorm.Expr(rankExpr, query, query))
	} else {
		db = db.Order("created_at DESC")
	}

	var count int64
	err := db.
		Where("product_state = ?", models.ProductStateActive).
		Count(&count).
		Error
	return count, err
}

// GetTopEndingSoons returns 5 products that are currently about to expire.
func (r *ProductRepository) GetTopEndingSoons(ctx context.Context) ([]models.Product, error) {
	return gorm.G[models.Product](r.DB).
		Preload("Seller", nil).
		Preload("Categories", nil).
		Preload("CurrentHighestBid.User", nil).
		Where("product_state = ?", models.ProductStateActive).
		Order("expired_at ASC").
		Limit(5).
		Find(ctx)
}

// GetMostActiveProducts returns 5 products that have the most number of bids.
func (r *ProductRepository) GetMostActiveProducts(ctx context.Context) ([]models.Product, error) {
	return gorm.G[models.Product](r.DB).
		Preload("Seller", nil).
		Preload("Categories", nil).
		Preload("CurrentHighestBid.User", nil).
		Where("product_state = ?", models.ProductStateActive).
		Order("bids_count DESC, expired_at ASC").
		Limit(5).
		Find(ctx)
}

// GetHighestBiddedProducts returns 5 products with the current highest active bid.
func (r *ProductRepository) GetHighestBiddedProducts(ctx context.Context) ([]models.Product, error) {
	var products []models.Product

	err := r.DB.WithContext(ctx).
		Joins("INNER JOIN bids ON products.current_highest_bid_id = bids.id").
		Preload("Seller").
		Preload("Categories").
		Preload("CurrentHighestBid.User").
		Where("product_state = ?", models.ProductStateActive).
		Order("bids.price DESC").
		Limit(5).
		Find(&products).
		Error

	return products, err
}

// GetProductByID gets a product's full details from an ID.
// This is a rather expensive query. Use sparingly, only when actually needed everything.
func (r *ProductRepository) GetProductByID(ctx context.Context, id int) (models.Product, error) {
	// This query should allow expired products to show up also.
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
		Preload("ChatSession", nil).
		Where("id = ?", id).
		First(ctx)
}

// GetSimilarProductsTo retrieves a list of products that are similar to another.
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
		Where("products.product_state = ?", models.ProductStateActive).
		Preload("Seller").
		Preload("Categories").
		Preload("CurrentHighestBid").
		Distinct("products.*").
		Limit(5).
		Find(&products).
		Error
	return products, err
}

// GetFavoriteProducts retrieves a list of products the user marked as favorite.
func (r *ProductRepository) GetFavoriteProducts(ctx context.Context, userID uint, limit int, offset int) ([]models.Product, error) {
	var products []models.Product
	err := r.DB.Model(&models.Product{}).WithContext(ctx).
		Joins("JOIN favorite_products ON products.id = favorite_products.product_id").
		Where("favorite_products.user_id = ?", userID).
		Where("product_state = ?", models.ProductStateActive).
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

// CountFavoriteProducts counts how many products the user has marked as favorite.
func (r *ProductRepository) CountFavoriteProducts(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.DB.Table("favorite_products").Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// ToggleFavoriteProduct toggles a favorite product state.
// This is tricky because it's a many-many table, not something you can just expr your way through.
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
	statement := gorm.G[models.Product](r.DB).
		Preload("Seller", nil).
		Preload("Categories", nil).
		Preload("CurrentHighestBid.User", nil).
		Where("product_state = ?", models.ProductStateActive).
		Where("seller_id = ?", userID)

	return statement.Order("expired_at").
		Limit(limit).
		Offset(offset).
		Find(ctx)
}

// CountRunningUserProducts counts the number of products posted by a user.
func (r *ProductRepository) CountRunningUserProducts(ctx context.Context, userID uint) (int64, error) {
	statement := gorm.G[models.Product](r.DB).
		Where("seller_id = ?", userID).
		Where("product_state = ?", models.ProductStateActive)
	return statement.Count(ctx, "id")
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
			Preload("CurrentHighestBid").
			Preload("Seller").
			Where("id = ?", productID).
			First(&product).
			Error
		if err != nil {
			return err
		}

		// If product is already inactive, remove
		if product.ProductState != models.ProductStateActive || product.ExpiredAt.Before(time.Now()) {
			return fmt.Errorf("product is already expired")
		}

		// TODO: Add this to config
		extensionSeconds := 300
		extensionDuration := time.Duration(extensionSeconds) * time.Second

		expiredAt := product.ExpiredAt
		expiryThreshold := 30 * time.Minute
		if product.AutoExtendsTime && time.Until(expiredAt) <= expiryThreshold {
			expiredAt = expiredAt.Add(extensionDuration)
		}

		// Check if it is the highest price.
		// Whether it makes sense (based on step bid is on the service, not this I think)
		if (product.CurrentHighestBid != nil && product.CurrentHighestBid.Price+product.StepBidValue > bidAmount) || bidAmount < product.StartingBid {
			return fmt.Errorf("bid is not high enough")
		}

		// Check if it's not the seller
		if product.SellerID == userID {
			return fmt.Errorf("you can't bid on your own auction")
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

		// Now mark the bids as counted, and extend it if it is needed.
		return tx.Model(&models.Product{Model: gorm.Model{ID: productID}}).
			Select("current_highest_bid_id", "bids_count", "expired_at").
			Updates(map[string]any{
				"current_highest_bid_id": bid.ID,
				"bids_count":             tx.Model(&models.Bid{}).Select("count(*)").Where("product_id = ?", productID),
				"expired_at":             expiredAt,
			}).
			Error
	})
}

// CreateAutomatedBid makes an automated bid.
func (r *ProductRepository) CreateAutomatedBid(
	ctx context.Context,
	productID uint,
	userID uint,
	maxAmount int64,
	lastBid *models.Bid,
	newBid *models.Bid,
	currentProduct *models.Product,
) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Lock Product and get current "Public" state
		product := models.Product{}
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("CurrentHighestBid").
			Preload("Seller").
			Where("id = ?", productID).
			First(&product).Error
		if err != nil {
			return err
		}

		// Validation
		if product.SellerID == userID {
			return fmt.Errorf("you can't bid on your own auction")
		}
		if time.Now().After(product.ExpiredAt) || product.ProductState != models.ProductStateActive {
			return fmt.Errorf("auction has already ended")
		}

		// 2. Upsert the User's Intent (Hidden Max)
		intent := models.BidIntent{
			ProductID: productID,
			UserID:    userID,
			BidAmount: maxAmount,
		}
		err = tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "product_id"}, {Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"bid_amount", "updated_at"}),
		}).Create(&intent).Error
		if err != nil {
			return err
		}

		// 3. Get Top 2 Intents to see the "Hidden" competition
		var topIntents []models.BidIntent
		err = tx.Where("product_id = ?", productID).
			Order("bid_amount DESC, created_at ASC").
			Limit(2).
			Find(&topIntents).Error
		if err != nil {
			return err
		}

		// 4. Determine the Winner and the new Public Price
		var finalWinnerID uint
		var finalPublicPrice int64

		highestIntent := topIntents[0]

		// The "Floor" is the current public price (manual or automated)
		currentPublicPrice := product.StartingBid
		if product.CurrentHighestBid != nil {
			currentPublicPrice = product.CurrentHighestBid.Price + product.StepBidValue
		}

		// Validation: The new intent MUST be at least currentPrice + step to be relevant
		// (Unless the user is just updating their max bid and is already the winner)
		isAlreadyWinner := product.CurrentHighestBid != nil && product.CurrentHighestBid.UserID == userID
		if !isAlreadyWinner && maxAmount < currentPublicPrice+product.StepBidValue {
			return fmt.Errorf("your max bid is too low to outbid the current price")
		}

		if len(topIntents) == 1 {
			// Scenario: Only one person has a proxy.
			// They win against the starting bid or the current manual bid.
			finalWinnerID = highestIntent.UserID
			finalPublicPrice = max(product.StartingBid, currentPublicPrice)
		} else {
			// Scenario: Fight between two proxies
			runnerUp := topIntents[1]
			finalWinnerID = highestIntent.UserID

			// The price is (Second Best Intent + Step),
			// but it must also be at least (Current Public Price + Step)
			calculatedPrice := max(runnerUp.BidAmount, currentPublicPrice+product.StepBidValue)

			// Cap it at the winner's maximum
			finalPublicPrice = min(highestIntent.BidAmount, calculatedPrice)
		}

		// 5. Create the History Record
		// If the winner hasn't changed AND the price hasn't changed (just updated max bid),
		// you might want to skip creating a new Bid row.
		if product.CurrentHighestBid != nil &&
			product.CurrentHighestBid.UserID == finalWinnerID &&
			product.CurrentHighestBid.Price == finalPublicPrice {
			// Just updating max bid ceiling, no public change
			*newBid = *product.CurrentHighestBid
		} else {
			bid := models.Bid{
				Price:     finalPublicPrice,
				UserID:    finalWinnerID,
				ProductID: product.ID,
			}
			if err := tx.Create(&bid).Error; err != nil {
				return err
			}
			*newBid = bid
		}

		// 6. Final Updates (Time & Product State)
		expiredAt := product.ExpiredAt
		if product.AutoExtendsTime && time.Until(expiredAt) <= 30*time.Minute {
			expiredAt = expiredAt.Add(300 * time.Second)
		}

		if product.CurrentHighestBid != nil {
			*lastBid = *product.CurrentHighestBid
		}
		*currentProduct = product

		return tx.Model(&models.Product{Model: gorm.Model{ID: productID}}).
			Select("current_highest_bid_id", "bids_count", "expired_at").
			Updates(map[string]any{
				"current_highest_bid_id": newBid.ID,
				"bids_count":             tx.Model(&models.Bid{}).Select("count(*)").Where("product_id = ?", productID),
				"expired_at":             expiredAt,
			}).Error
	})
}

func (r *ProductRepository) ClearAllBids(ctx context.Context, productID uint) (int64, error) {
	db := r.DB.WithContext(ctx).
		Model(&models.Bid{}).
		Where("product_id = ?", productID).
		Delete(&models.Bid{})
	return db.RowsAffected, db.Error
}

func (r *ProductRepository) CreateBINPurchase(ctx context.Context, productID uint, userID uint) error {
	now := time.Now()
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Retrieve the product about to update
		product := models.Product{}
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&models.Product{}).
			Preload("CurrentHighestBid").
			Preload("Seller").
			Where("id = ?", productID).
			First(&product).
			Error
		if err != nil {
			return err
		}

		if product.BINPrice == nil {
			return fmt.Errorf("product %d does not have a bin price", productID)
		}

		// Check if it's not the seller
		if product.SellerID == userID {
			return fmt.Errorf("you can't bid on your own auction")
		}

		// Now try to insert into the bids table I guess.
		bid := models.Bid{Price: *product.BINPrice, UserID: userID, ProductID: product.ID}
		err = tx.Model(&bid).Create(&bid).Error
		if err != nil {
			return err
		}

		// Now mark the bids as counted, and extend it if it is needed.
		return tx.Model(&models.Product{Model: gorm.Model{ID: productID}}).
			Select("current_highest_bid_id", "bids_count", "expired_at").
			Updates(map[string]any{
				"current_highest_bid_id": bid.ID,
				"bids_count":             gorm.Expr("bids_count + 1"),
				"expired_at":             now,
			}).
			Error
	})
}

// GetMyBids retrieves a user's bids.
// Should this allow the user to see won bids? Prob not, let's call that transactions.
func (r *ProductRepository) GetMyBids(ctx context.Context, userID uint, ended bool, limit int, offset int) ([]models.Product, error) {
	var products []models.Product
	db := r.DB.WithContext(ctx).
		Model(&models.Product{}).
		Select("DISTINCT ON (products.id) products.*").
		Preload("Seller").
		Preload("CurrentHighestBid.User").
		Preload("Bids.User").
		Joins("JOIN bids ON bids.product_id = products.id AND bids.user_id = ? AND bids.deleted_at IS NULL", userID).
		Order("products.id, products.expired_at ASC")

	if ended {
		db = db.Where("products.product_state in ?", []models.ProductState{models.ProductStateEnded, models.ProductStateExpired})
	} else {
		db = db.Where("products.product_state = ?", models.ProductStateActive)
	}

	err := db.
		Limit(limit).
		Offset(offset).
		Find(&products).
		Error
	return products, err
}

func (r *ProductRepository) CountMyBids(ctx context.Context, userID uint, ended bool) (int64, error) {
	var count int64
	db := r.DB.WithContext(ctx).
		Select("DISTINCT ON (products.id) products.*").
		Model(&models.Product{}).
		Joins("JOIN bids ON bids.product_id = products.id AND bids.user_id = ? AND bids.deleted_at IS NULL", userID)

	if ended {
		db = db.Where("products.product_state in ?", []models.ProductState{models.ProductStateEnded, models.ProductStateExpired})
	} else {
		db = db.Where("products.product_state = ?", models.ProductStateActive)
	}

	err := db.
		Count(&count).
		Error
	return count, err
}

// SetProductSentEmail marks a product's sent-email field as true.
func (r *ProductRepository) SetProductSentEmail(ctx context.Context, productID uint) (int, error) {
	db := r.DB.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ?", productID).
		Update("email_sent", true)
	return int(db.RowsAffected), db.Error
}

// GetAllExpiredProducts retrieves all products that are expired or ended
// but not have an email sent yet.
func (r *ProductRepository) GetAllExpiredProducts(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	err := r.DB.WithContext(ctx).
		Model(&models.Product{}).
		Preload("CurrentHighestBid.User").
		Preload("Seller").
		Where("product_state <> ?", models.ProductStateActive).
		Where("email_sent = ?", false).
		Find(&products).
		Error
	return products, err
}

// UpdateAllExpiredProducts updates all products' state to match their status.
func (r *ProductRepository) UpdateAllExpiredProducts(ctx context.Context) error {
	now := time.Now()

	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 入札なし → EXPIRED
		if err := tx.
			Model(&models.Product{}).
			Where("product_state = ?", models.ProductStateActive).
			Where("expired_at < ?", now).
			Where("current_highest_bid_id IS NULL").
			Update("product_state", models.ProductStateExpired).
			Error; err != nil {
			return err
		}

		// 入札あり → ENDED
		if err := tx.
			Model(&models.Product{}).
			Where("product_state = ?", models.ProductStateActive).
			Where("expired_at < ?", now).
			Where("current_highest_bid_id IS NOT NULL").
			Update("product_state", models.ProductStateEnded).
			Error; err != nil {
			return err
		}

		return nil
	})
}

// DenyBidder marks a bidder denied from the product.
func (r *ProductRepository) DenyBidder(ctx context.Context, productID uint, userID uint) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		tx.Model(&models.DeniedBidder{}).
			Where("product_id = ? AND user_id = ?", productID, userID).
			Count(&count)

		if count > 0 {
			return fmt.Errorf("user %d is already denied for product %d", userID, productID)
		}

		err := tx.Model(&models.Product{Model: gorm.Model{ID: productID}}).
			Association("DeniedBidders").
			Append(&models.DeniedBidder{UserID: userID})
		if err != nil {
			return err
		}

		if err := tx.Where("product_id = ? AND user_id = ?", productID, userID).Delete(&models.Bid{}).Error; err != nil {
			return err
		}

		return tx.Model(&models.Product{}).
			Where("id = ?", productID).
			Select("bids_count", "current_highest_bid_id").
			Updates(map[string]any{
				"bids_count": tx.Model(&models.Bid{}).Select("count(id)").Where("product_id = ?", productID),
				"current_highest_bid_id": tx.Model(&models.Bid{}).
					Select("id").
					Where("product_id = ?", productID).
					Order("price DESC, created_at ASC").
					Limit(1),
			}).Error
	})
}

func (r *ProductRepository) GetUserProducts(ctx context.Context, userID uint, state models.ProductState, limit int, offset int) ([]models.Product, error) {
	var products []models.Product
	err := r.DB.
		WithContext(ctx).
		Model(&models.Product{}).
		Preload("Seller").
		Preload("CurrentHighestBid.User").
		Preload("Bids.User").
		Where("seller_id = ?", userID).
		Where("product_state = ?", state).
		Order("expired_at DESC, created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).
		Error
	return products, err
}

func (r *ProductRepository) CountUserProducts(ctx context.Context, userID uint, state models.ProductState) (int64, error) {
	var count int64
	err := r.DB.
		WithContext(ctx).
		Model(&models.Product{}).
		Where("seller_id = ?", userID).
		Where("product_state = ?", state).
		Count(&count).
		Error
	return count, err
}

// FinalizeProduct marks a product as finalized as of now.
func (r *ProductRepository) FinalizeProduct(ctx context.Context, id uint) (int64, error) {
	db := r.DB.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ?", id).
		Update("finalized_at", time.Now())
	return db.RowsAffected, db.Error
}
