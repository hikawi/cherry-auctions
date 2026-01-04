package products

import "github.com/gin-gonic/gin"

// PostBid godoc
//
// @summary Creates a bid on a product
// @description Bids on a product, hopefully with some race-condition management.
// @router /v1/products/{id}/bids [POST]
func (h *ProductsHandler) PostBid(g *gin.Context) {
}
