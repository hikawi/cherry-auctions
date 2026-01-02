package products

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"luny.dev/cherryauctions/internal/logging"
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/routes/shared"
	"luny.dev/cherryauctions/internal/services"
	"luny.dev/cherryauctions/pkg/closer"
	"luny.dev/cherryauctions/pkg/ranges"
)

// GetProducts godoc
//
//	@summary		Queries products using fuzzy matchers and full-text matchers.
//	@description	Queries from a list of products using a set of keywords, using Full-text Queries or Fuzzy and Similarity queries.
//	@tags			products
//	@produce		json
//	@param			query		query		string	false	"Search Query"
//	@param			page		query		int		false	"Page Number"
//	@param			per_page	query		int		false	"Items per Page"
//	@success		200			{object}	products.GetProductsResponse
//	@failure		500			{object}	shared.ErrorResponse	"The server could not make the request"
//	@router			/products [get]
func (h *ProductsHandler) GetProducts(g *gin.Context) {
	ctx := g.Request.Context()
	query := GetProductsQuery{
		Page:    1,
		PerPage: 20,
	}

	if err := g.ShouldBindQuery(&query); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "invalid query"})
		return
	}

	products, err := h.ProductRepo.SearchProducts(ctx, query.Query, query.PerPage, (query.Page-1)*query.PerPage)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "couldn't query the database"})
		return
	}

	count, err := h.ProductRepo.CountProductsWithQuery(ctx, query.Query)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "unable to count products"})
		return
	}

	productsDto := make([]ProductDTO, 0)
	for _, prod := range products {
		productsDto = append(productsDto, ToProductDTO(&prod))
	}

	response := GetProductsResponse{
		Data:       productsDto,
		Total:      count,
		TotalPages: int(math.Ceil(float64(count) / float64(query.PerPage))),
		Page:       query.Page,
		PerPage:    query.PerPage,
	}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "response": response, "query": query})
	g.JSON(http.StatusOK, response)
}

// GetProductsTop godoc
//
//	@summary		Retrieves a list of products of top-categories.
//	@description	Queriee top-products via some metrics, instead of a pure query like /products.
//	@tags			products
//	@produce		json
//	@success		200	{object}	products.GetTopProductsResponse
//	@failure		500	{object}	shared.ErrorResponse	"The server could not make the request"
//	@router			/products/top [get]
func (h *ProductsHandler) GetProductsTop(g *gin.Context) {
	ctx := g.Request.Context()

	highestBids, err := h.ProductRepo.GetHighestBiddedProducts(ctx)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusInternalServerError})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "unable to read products for top bidded"})
		return
	}

	endingSoon, err := h.ProductRepo.GetTopEndingSoons(ctx)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusInternalServerError})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "unable to read products for ending soon"})
		return
	}

	topBids, err := h.ProductRepo.GetMostActiveProducts(ctx)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusInternalServerError})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "unable to read products for ending soon"})
		return
	}

	topBidsPtr := ranges.Each(topBids, func(m models.Product) *models.Product { return &m })
	endingSoonsPtr := ranges.Each(endingSoon, func(m models.Product) *models.Product { return &m })
	highestBidsPtr := ranges.Each(highestBids, func(m models.Product) *models.Product { return &m })

	if claims, ok := g.Get("claims"); ok {
		jwt := claims.(*services.JWTSubject)
		var allProducts []*models.Product
		allProducts = append(allProducts, topBidsPtr...)
		allProducts = append(allProducts, endingSoonsPtr...)
		allProducts = append(allProducts, highestBidsPtr...)
		h.ProductRepo.AttachFavoriteStatus(ctx, jwt.UserID, allProducts...)
	}

	response := GetTopProductsResponse{
		TopBidded:   ToProductDTOs(topBidsPtr),
		EndingSoon:  ToProductDTOs(endingSoonsPtr),
		HighestBids: ToProductDTOs(highestBidsPtr),
	}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "response": response})
	g.JSON(http.StatusOK, response)
}

// GetFavoriteProducts godoc
//
//	@summary		Retrieves my favorite products.
//	@description	Queries on my favorite products,
//	@tags			products
//	@security		ApiKeyAuth
//	@produce		json
//	@param			query		query		string	false	"Search Query"
//	@param			page		query		int		false	"Page Number"
//	@param			per_page	query		int		false	"Items per Page"
//	@success		200			{object}	products.GetProductsResponse
//	@failure		401			{object}	shared.ErrorResponse	"When unauthenticated"
//	@failure		500			{object}	shared.ErrorResponse	"The server could not make the request"
//	@router			/products/favorite [get]
func (h *ProductsHandler) GetFavoriteProducts(g *gin.Context) {
	ctx := g.Request.Context()
	claimsAny, _ := g.Get("claims")
	claims := claimsAny.(*services.JWTSubject)
	query := GetProductsQuery{
		Page:    1,
		PerPage: 20,
	}

	if err := g.ShouldBindQuery(&query); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "invalid query"})
		return
	}

	products, err := h.ProductRepo.GetFavoriteProducts(ctx, claims.UserID, query.PerPage, (query.Page-1)*query.PerPage)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "couldn't read favorite products"})
		return
	}

	count, err := h.ProductRepo.CountFavoriteProducts(ctx, claims.UserID)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "unable to count products"})
		return
	}

	// Attach the favorite issues=
	pointers := ranges.Each(products, func(p models.Product) *models.Product { return &p })
	h.ProductRepo.AttachFavoriteStatus(ctx, claims.UserID, pointers...)

	productsDto := make([]ProductDTO, 0)
	for _, prod := range pointers {
		productsDto = append(productsDto, ToProductDTO(prod))
	}

	response := GetProductsResponse{
		Data:       productsDto,
		Total:      count,
		TotalPages: int(math.Ceil(float64(count) / float64(query.PerPage))),
		Page:       query.Page,
		PerPage:    query.PerPage,
	}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "response": response, "query": query})
	g.JSON(http.StatusOK, response)
}

// PostFavoriteProduct godoc
//
//	@summary		Toggle a product as favorite.
//	@description	Mark a product as my favorite, if it is not a favorite, otherwise, unmark it.
//	@tags			products
//	@security		ApiKeyAuth
//	@produce		json
//	@success		204	{object}	shared.MessageResponse
//	@failure		401	{object}	shared.ErrorResponse	"When unauthenticated"
//	@failure		404	{object}	shared.ErrorResponse	"Product could not be found"
//	@failure		500	{object}	shared.ErrorResponse	"The server could not make the request"
//	@router			/products/favorite [post]
func (h *ProductsHandler) PostFavoriteProduct(g *gin.Context) {
	ctx := g.Request.Context()
	claimsAny, _ := g.Get("claims")
	claims := claimsAny.(*services.JWTSubject)

	var body PostProductIDBody
	if err := g.ShouldBind(&body); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "body": body})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "invalid product id"})
		return
	}

	err := h.ProductRepo.ToggleFavoriteProduct(ctx, claims.UserID, body.ID)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "body": body})
		g.AbortWithStatusJSON(http.StatusNotFound, shared.ErrorResponse{Error: "invalid product id, does not exist"})
		return
	}

	response := shared.MessageResponse{Message: "toggled favorite"}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "body": body, "response": response})
	g.JSON(http.StatusOK, response)
}

// GetProductID godoc
//
//	@summary		Retrieves details on a single product.
//	@description	Queries everything about a single product, with full joins.
//	@tags			products
//	@produce		json
//	@param			id	path		int	true	"Product ID to view"
//	@success		200	{object}	products.GetProductDetailsResponse
//	@failure		404	{object}	shared.ErrorResponse	"The server couldn't find the requested product"
//	@failure		500	{object}	shared.ErrorResponse	"The server could not make the request"
//	@router			/products/{id} [get]
func (h *ProductsHandler) GetProductID(g *gin.Context) {
	ctx := g.Request.Context()

	paramId := g.Param("id")
	id, err := strconv.ParseInt(paramId, 10, 0)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "id": paramId})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	product, err := h.ProductRepo.GetProductByID(ctx, int(id))
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusNotFound, "error": err.Error(), "id": paramId})
		g.AbortWithStatusJSON(http.StatusNotFound, shared.ErrorResponse{Error: "can't find product with that id"})
		return
	}

	similars, err := h.ProductRepo.GetSimilarProductsTo(ctx, &product)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "id": paramId})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "can't find similar products"})
		return
	}

	similarsPtr := ranges.Each(similars, func(m models.Product) *models.Product { return &m })

	if jwt, ok := g.Get("claims"); ok {
		claims := jwt.(*services.JWTSubject)

		var all []*models.Product
		all = append(all, &product)
		all = append(all, similarsPtr...)
		h.ProductRepo.AttachFavoriteStatus(ctx, claims.UserID, all...)
	}

	response := GetProductDetailsResponse{
		ProductDTO:      ToProductDTO(&product),
		ProductImages:   ToProductImageDTOs(product.ProductImages),
		Questions:       ToQuestionDTOs(product.Questions),
		Bids:            ToBidDTOs(product.Bids),
		SimilarProducts: ToProductDTOs(similarsPtr),
	}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "response": similars})
	g.JSON(http.StatusOK, response)
}

// Function courtesy of Gemini, based on the function in users/handler.go.
func (h *ProductsHandler) uploadImages(ctx context.Context, body PostProductBody) ([]string, error) {
	// Create an errgroup with the request context
	g, gCtx := errgroup.WithContext(ctx)

	// Slice to hold the resulting URLs.
	// We use a fixed length to avoid race conditions when writing to indices.
	imageURLs := make([]string, len(body.ProductImages))

	for i, fileHeader := range body.ProductImages {
		g.Go(func() error {
			// 1. Open File
			file, err := fileHeader.Open()
			if err != nil {
				return fmt.Errorf("failed to open file %d: %w", i, err)
			}
			defer closer.CloseResources(file)

			// 2. Load into VIPS
			img, err := vips.NewImageFromReader(file)
			if err != nil {
				return fmt.Errorf("invalid image format for file %d: %w", i, err)
			}
			defer img.Close()

			// 3. Process (Thumbnail)
			// For product images, you might want a larger size than avatars
			err = img.Thumbnail(800, 800, vips.InterestingAttention)
			if err != nil {
				return fmt.Errorf("processing failed for file %d: %w", i, err)
			}

			// 4. Export to WebP
			ep := vips.NewWebpExportParams()
			ep.Quality = 75
			ep.StripMetadata = true
			webpBuffer, _, err := img.ExportWebp(ep)
			if err != nil {
				return fmt.Errorf("encoding failed for file %d: %w", i, err)
			}

			// 5. Upload to S3
			hash := sha256.Sum256(webpBuffer)
			hashString := hex.EncodeToString(hash[:])
			key := fmt.Sprintf("products/%s.webp", hashString)

			// Use gCtx to respect group cancellation if another upload fails
			err = h.S3Service.PutObject(gCtx, key, bytes.NewReader(webpBuffer))
			if err != nil {
				return fmt.Errorf("storage upload failed for file %d: %w", i, err)
			}

			// 6. Store the URL in the slice at the correct index
			imageURLs[i] = fmt.Sprintf("%s/%s", h.S3PermURL, key)
			return nil
		})
	}

	// Wait for all uploads to finish or the first error to occur
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return imageURLs, nil
}

// PostProduct godoc
//
//	@summary		Posts a new product.
//	@description	Posts a new product up for auctions.
//	@tags			products
//	@security ApiKeyAuth
//	@accept multipart/form-data
//	@produce		json
//	@success		201	{object}	shared.MessageResponse "Successfully created an auction"
//	@failure 400 {object} shared.ErrorResponse "When the multipart data is invalid"
//	@failure 401 {object} shared.ErrorResponse "When the user is unauthorized"
//	@failure 403 {object} shared.ErrorResponse "When the user isn't subscribed"
//	@failure		404	{object}	shared.ErrorResponse	"The server couldn't find the requested product"
//	@failure		500	{object}	shared.ErrorResponse	"The server could not make the request"
//	@router			/products/{id} [get]
func (h *ProductsHandler) PostProduct(g *gin.Context) {
	ctx := g.Request.Context()
	claimsAny, _ := g.Get("claims")
	claims := claimsAny.(*services.JWTSubject)

	// Make sure the user has the permission.
	if claims.SubscriptionExpiredAt == nil || claims.SubscriptionExpiredAt.Before(time.Now()) {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": "user can't post"})
		g.AbortWithStatusJSON(http.StatusForbidden, shared.ErrorResponse{Error: "you can't post"})
		return
	}

	var body PostProductBody
	if err := g.ShouldBind(&body); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	// Validation for product images size.
	for _, img := range body.ProductImages {
		if img.Size > (10 << 20) /* 10MB */ {
			logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": "image too large", "status": http.StatusRequestEntityTooLarge, "size": img.Size})
			g.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, shared.ErrorResponse{Error: "max 10MB allowed"})
			return
		}
	}

	// Try to upload the images separately.
	urls, err := h.uploadImages(ctx, body)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusInternalServerError, "body": body})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "failed to upload image"})
		return
	}

	product := models.Product{
		Name:                body.Name,
		Description:         body.Description,
		StartingBid:         body.StartingBid,
		StepBidType:         body.StepBidType,
		StepBidValue:        body.StepBidValue,
		BINPrice:            body.BINPrice,
		AllowsUnratedBuyers: body.AllowsUnrated,
		AutoExtendsTime:     body.AutoExtends,
		ExpiredAt:           body.ExpiredAt,
		SellerID:            claims.UserID,
		ThumbnailURL:        urls[0],
		ProductImages: ranges.Each(urls[1:], func(url string) models.ProductImage {
			return models.ProductImage{
				URL:     url,
				AltText: body.Name,
			}
		}),
	}

	// Create the product
	err = h.ProductRepo.CreateProduct(ctx, &product)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusInternalServerError, "body": body, "urls": urls})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "failed to create product"})
		return
	}

	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusCreated, "body": body, "response": product})
	g.JSON(http.StatusCreated, product)
}
