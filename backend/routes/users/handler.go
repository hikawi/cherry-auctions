package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/repositories"
	"luny.dev/cherryauctions/routes/shared"
	"luny.dev/cherryauctions/services"
	"luny.dev/cherryauctions/utils"
)

// GetMe retrieves your own profile if logged in.
//
//	@summary		Gets your own profile.
//	@description	Retrieves information about your own profile if authenticated.
//	@tags			users
//	@produce		json
//	@success		200	{object}	users.GetMeResponse
//	@failure		401	{object}	shared.ErrorResponse	"When unauthenticated"
//	@failure		422	{object}	shared.ErrorResponse	"When your info had an invalid state on the server"
//	@failure		500	{object}	shared.ErrorResponse	"The request could not be completed due to server faults"
//	@router			/users/me [GET]
func (h *UsersHandler) GetMe(g *gin.Context) {
	claimsAny, ok := g.Get("claims")
	if !ok {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusUnauthorized, "err": "unauthenticated read to get me"})
		g.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse{Error: "unauthenticated"})
		return
	}

	claims := claimsAny.(*services.JWTSubject)
	userRepo := repositories.UserRepository{DB: h.DB}

	user, err := userRepo.GetUserByID(claims.ID)
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusUnprocessableEntity, "err": err.Error(), "claims": claims})
		g.AbortWithStatusJSON(http.StatusUnprocessableEntity, shared.ErrorResponse{Error: "unknown user but authenticated"})
		return
	}

	res := GetMeResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		OauthType: user.OauthType,
		Verified:  user.Verified,
	}
	utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusOK, "response": res})
	g.JSON(200, res)
}
