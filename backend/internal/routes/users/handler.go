package users

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

// GetMe retrieves your own profile if logged in.
//
//	@summary		Gets your own profile.
//	@description	Retrieves information about your own profile if authenticated.
//	@tags			users
//	@produce		json
//	@security		ApiKeyAuth
//	@success		200	{object}	users.UserDTO
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

	response := ToUserDTO(&user)
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "response": response})
	g.JSON(http.StatusOK, response)
}

// PostRequest godoc
//
//	@summary		Requests seller privileges
//	@description	Sends a request to the admin to approve or deny seller privileges.
//	@tags			users
//	@produce		json
//	@security		ApiKeyAuth
//	@success		204	{object}	shared.MessageResponse	"When success"
//	@failure		401	{object}	shared.ErrorResponse	"When unauthenticated"
//	@failure		500	{object}	shared.ErrorResponse	"The request could not be completed due to server faults"
//	@router			/users/request [POST]
func (h *UsersHandler) PostRequest(g *gin.Context) {
	claimsAny, _ := g.Get("claims")
	claims := claimsAny.(*services.JWTSubject)
	ctx := g.Request.Context()

	rows, err := h.UserRepo.RequestUserApproval(ctx, claims.UserID)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusInternalServerError})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "can't mark user as requesting"})
		return
	}

	if rows == 0 {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": "no rows are written", "status": http.StatusInternalServerError})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "no rows in db were changed"})
		return
	}

	response := shared.MessageResponse{Message: "requested privileges successfully"}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusNoContent, "response": response})
	g.JSON(http.StatusNoContent, response)
}

// PostApprove godoc
//
//	@summary		Approves seller privileges
//	@description	Approves seller privileges
//	@tags			users
//	@produce		json
//	@security		ApiKeyAuth
//	@success		204	{object}	shared.MessageResponse	"When success"
//	@failure		400	{object}	shared.ErrorResponse	"Invalid request"
//	@failure		401	{object}	shared.ErrorResponse	"When unauthenticated"
//	@failure		500	{object}	shared.ErrorResponse	"The request could not be completed due to server faults"
//	@router			/users/approve [POST]
func (h *UsersHandler) PostApprove(g *gin.Context) {
	ctx := g.Request.Context()

	var body PostApproveRequest
	if err := g.ShouldBindBodyWithJSON(&body); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusBadRequest, "body": body})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	err := h.UserRepo.ApproveUser(ctx, uint(body.ID))
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusBadRequest, "body": body})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "couldn't add a new subscription"})
		return
	}

	response := shared.MessageResponse{Message: "approved successfully"}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusNoContent, "response": response})
	g.JSON(http.StatusNoContent, response)
}

// GetUsers godoc
//
//	@summary		Retrieves all users
//	@description	Retrieves all users.
//	@tags			users
//	@produce		json
//	@security		ApiKeyAuth
//	@success		200	{object}	users.GetUsersResponse	"All users"
//	@failure		401	{object}	shared.ErrorResponse	"When unauthorized"
//	@failure		500	{object}	shared.ErrorResponse	"The request could not be completed due to server faults"
//	@router			/users [GET]
func (h *UsersHandler) GetUsers(g *gin.Context) {
	ctx := g.Request.Context()
	query := GetUsersQuery{
		Page:    1,
		PerPage: 20,
	}

	if err := g.ShouldBindQuery(&query); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "invalid query"})
		return
	}

	users, err := h.UserRepo.GetUsers(ctx, query.PerPage, (query.Page-1)*query.PerPage)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "couldn't query the database"})
		return
	}

	count, err := h.UserRepo.CountUsers(ctx)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "unable to count products"})
		return
	}

	response := GetUsersResponse{
		Data:       ranges.Each(users, func(m models.User) UserDTO { return ToUserDTO(&m) }),
		Total:      count,
		TotalPages: int(math.Ceil(float64(count) / float64(query.PerPage))),
		Page:       query.Page,
		PerPage:    query.PerPage,
	}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "response": response, "query": query})
	g.JSON(http.StatusOK, response)
}
