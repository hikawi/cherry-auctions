package shared

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/services"
)

// AuthenticatedRoute forces the Gin Context to have an Authorization header to work.
func AuthenticatedRoute(g *gin.Context) {
	authHeader := strings.Split(g.GetHeader("Authorization"), " ")
	if len(authHeader) != 2 || authHeader[0] != "Bearer" {
		g.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Error: "no bearer authentication"})
		return
	}

	claims, err := services.VerifyJWT(authHeader[1])
	if err != nil {
		g.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid access token"})
		return
	}

	g.Set("claims", claims)
	g.Next()
}
