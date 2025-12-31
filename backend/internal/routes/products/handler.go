package products

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/logging"
	"luny.dev/cherryauctions/internal/routes/shared"
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
//	@success		200			{object}	products.GetProductsResponse
//	@failure		500			{object}	shared.ErrorResponse	"The server could not make the request"
//	@router			/products/top [get]
func (h *ProductsHandler) GetProductsTop(g *gin.Context) {
	ctx := g.Request.Context()

	topBidded, err := h.ProductRepo.GetHighestBiddedProducts(ctx)
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

	response := gin.H{
		"highest_bids": ToProductDTOs(topBidded),
		"ending_soon":  ToProductDTOs(endingSoon),
		"top_bids":     ToProductDTOs(topBids),
	}
	g.JSON(http.StatusOK, response)
}
