package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/logging"
	"luny.dev/cherryauctions/internal/routes/shared"
	"luny.dev/cherryauctions/internal/services"
)

// GetMe retrieves your own profile if logged in.
//
//	@summary		Gets your own profile.
//	@description	Retrieves information about your own profile if authenticated.
//	@tags			users
//	@produce		json
//	@security		ApiKeyAuth
//	@success		200	{object}	users.GetMeResponse
//	@failure		401	{object}	shared.ErrorResponse	"When unauthenticated"
//	@failure		422	{object}	shared.ErrorResponse	"When your info had an invalid state on the server"
//	@failure		500	{object}	shared.ErrorResponse	"The request could not be completed due to server faults"
//	@router			/users/me [GET]
func (h *UsersHandler) GetMe(g *gin.Context) {
	claimsAny, _ := g.Get("claims")
	claims := claimsAny.(*services.JWTSubject)
	ctx := g.Request.Context()

	user, err := h.UserRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusUnprocessableEntity})
		g.AbortWithStatusJSON(http.StatusUnprocessableEntity, shared.ErrorResponse{Error: "unknown user but authenticated"})
		return
	}

	roles := make([]string, 0)
	for _, role := range user.Roles {
		roles = append(roles, role.ID)
	}

	res := GetMeResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     *user.Email,
		Roles:     roles,
		OauthType: user.OauthType,
		Verified:  user.Verified,
	}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "response": res})
	g.JSON(http.StatusOK, res)
}
