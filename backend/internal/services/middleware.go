package services

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/logging"
	"luny.dev/cherryauctions/internal/routes/shared"
)

type MiddlewareService struct {
	JWTService *JWTService
}

func (s *MiddlewareService) parseAuthHeaders(g *gin.Context) (*JWTSubject, error) {
	authHeader := strings.Split(g.GetHeader("Authorization"), " ")
	if len(authHeader) != 2 || authHeader[0] != "Bearer" {
		return nil, errors.New("no bearer authentication header found")
	}

	claims, err := s.JWTService.VerifyJWT(authHeader[1])
	if err != nil {
		return nil, errors.New("invalid access token")
	}

	return claims, nil
}

// AuthorizedRoleRoute creates a GIN handler that forces the user to have a specified permission name
// to be allowed to proceed. Otherwise, block with `403 forbidden`.
//
// If role is empty, this just wants the user to be authenticated.
func (s *MiddlewareService) AuthorizedRoute(role string) func(*gin.Context) {
	return func(g *gin.Context) {
		claims, err := s.parseAuthHeaders(g)
		if err != nil {
			logging.LogMessage(g, logging.LOG_DEBUG, gin.H{"error": err.Error()})
			g.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse{Error: err.Error()})
			return
		}

		// Doesn't support wildcard permissions, but I don't care.
		g.Set("claims", claims)
		roles := strings.Split(claims.Roles, " ")
		if role != "" {
			for _, hasRole := range roles {
				if strings.EqualFold(hasRole, role) {
					g.Next()
					return
				}
			}

			// Not allowed I guess.
			logging.LogMessage(g, logging.LOG_DEBUG, gin.H{"error": "not enough permissions", "roles": roles})
			g.AbortWithStatusJSON(http.StatusForbidden, shared.ErrorResponse{Error: "not enough permissions"})
		}

		g.Next()
	}
}
