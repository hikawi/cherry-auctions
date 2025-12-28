// Package routes provides the up-to-date version of routes for GIN to hook onto.
package routes

import "github.com/gin-gonic/gin"

// GetHealth godoc
//
//	@summary	Checks health of the server.
//	@tags		others
//	@produce	json
//	@success	200	{object}	shared.MessageResponse	"Always"
//	@router		/health [GET]
func GetHealth(g *gin.Context) {
	g.JSON(200, gin.H{"message": "healthy"})
}
