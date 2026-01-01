package products

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/logging"
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/routes/shared"
	"luny.dev/cherryauctions/internal/services"
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
		fmt.Println(prod.IsFavorite)
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
