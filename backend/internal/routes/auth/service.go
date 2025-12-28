package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"luny.dev/cherryauctions/internal/logging"
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/routes/shared"
)

func (h *AuthHandler) assignJWTKeyPair(g *gin.Context, db *gorm.DB, loggingBody any, id uint, email, roles string) {
	// Generate a JWT key pair.
	accessToken, err := h.JWTService.SignJWT(id, email, roles)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": "server can't sign jwt", "body": loggingBody})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server can't sign jwt"})
		return
	}

	refreshToken, err := h.RandomService.GenerateSecretKey(64)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "body": loggingBody})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server can't generate jwt key pair"})
		return
	}

	// Save the refresh token.
	hashedToken := sha256.Sum256(refreshToken)
	_, err = h.RefreshTokenRepo.SaveUserToken(g.Request.Context(), id, base64.URLEncoding.EncodeToString(hashedToken[:]))
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "body": loggingBody})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server can't hash refresh token"})
		return
	}

	g.SetCookieData(&http.Cookie{
		Name:     "RefreshToken",
		Value:    base64.URLEncoding.EncodeToString(refreshToken),
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 30 * 3),
		Domain:   h.Domain,
		Secure:   h.CookieSecure,
		SameSite: http.SameSiteNoneMode,
	})
	g.JSON(http.StatusOK, LoginResponse{AccessToken: accessToken})
}

func (h *AuthHandler) toRoleString(roles []models.Role) string {
	names := make([]string, 0)
	for _, role := range roles {
		names = append(names, role.ID)
	}
	return strings.Join(names, " ")
}
