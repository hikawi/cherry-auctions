package products

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/logging"
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/routes/shared"
	"luny.dev/cherryauctions/internal/services"
	"luny.dev/cherryauctions/pkg/ranges"
)

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
