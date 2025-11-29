// Package routes provides the up-to-date version of routes for GIN to hook onto.
package routes

import "github.com/gin-gonic/gin"

func GetHealth(g *gin.Context) {
	g.JSON(200, gin.H{"message": "healthy"})
}
