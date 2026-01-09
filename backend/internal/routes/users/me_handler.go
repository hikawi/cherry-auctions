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

// PutProfile godoc
//
//	@summary		Updates your user profile
//	@description	Updates the current authenticated user's profile
//	@tags			users
//	@accept			json
//	@produce		json
//	@security		ApiKeyAuth
//	@param			profile	body		users.PostProfileRequest	true	"Profile data"
//	@success		200		{object}	shared.MessageResponse		"When successfully changed"
//	@failure		400		{object}	shared.ErrorResponse		"Invalid body"
//	@failure		401		{object}	shared.ErrorResponse		"When unauthorized"
//	@failure		500		{object}	shared.ErrorResponse		"The request could not be completed due to server faults"
//	@router			/users/me [PUT]
func (h *UsersHandler) PutProfile(g *gin.Context) {
	ctx := g.Request.Context()
	claimsAny, _ := g.Get("claims")
	claims := claimsAny.(*services.JWTSubject)

	var body PostProfileRequest
	if err := g.ShouldBind(&body); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	_, err := h.UserRepo.UpdateProfile(ctx, claims.UserID, body.Name, body.Address)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusInternalServerError})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "couldn't update profile"})
		return
	}

	response := shared.MessageResponse{Message: "updated profile"}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "response": response})
	g.JSON(http.StatusOK, response)
}

// PutPassword godoc
//
//	@summary		Updates your profile password
//	@description	Updates the current authenticated user's profile password
//	@tags			users
//	@accept			json
//	@produce		json
//	@security		ApiKeyAuth
//	@param			profile	body		users.PutPasswordRequest	true	"New password data"
//	@success		200		{object}	shared.MessageResponse		"When successfully changed"
//	@failure		400		{object}	shared.ErrorResponse		"Invalid body"
//	@failure		401		{object}	shared.ErrorResponse		"When unauthorized"
//	@failure		403		{object}	shared.ErrorResponse		"Incorrect current password"
//	@failure		421		{object}	shared.ErrorResponse		"Account uses OAuth"
//	@failure		500		{object}	shared.ErrorResponse		"The request could not be completed due to server faults"
//	@router			/users/me/password [PUT]
func (h *UsersHandler) PutPassword(g *gin.Context) {
	ctx := g.Request.Context()
	claimsAny, _ := g.Get("claims")
	claims := claimsAny.(*services.JWTSubject)

	var body PutPasswordRequest
	if err := g.ShouldBind(&body); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	user, err := h.UserRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusInternalServerError})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "failed to query for user"})
		return
	}

	if user.Password == nil || user.OauthType != "none" {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusMisdirectedRequest, "error": "user does not use password"})
		g.AbortWithStatusJSON(http.StatusMisdirectedRequest, shared.ErrorResponse{Error: "user does not use password"})
		return
	}

	ok, err := h.PasswordService.VerifyPassword(*user.Password, body.CurrentPassword)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "failed to verify user password"})
		return
	}

	if !ok {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": "wrong password"})
		g.AbortWithStatusJSON(http.StatusForbidden, shared.ErrorResponse{Error: "wrong password"})
		return
	}

	hash, err := h.PasswordService.HashPassword(body.NewPassword)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": "failed to hash password"})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "failed to hash password"})
		return
	}

	_, err = h.UserRepo.UpdatePassword(ctx, user.ID, hash)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "failed to update password"})
		return
	}

	response := shared.MessageResponse{Message: "updated password"}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "response": response})
	g.JSON(http.StatusOK, response)
}

// GetMyProducts godoc
//
//	@summary		Retrieves my products
//	@description	Retrieves my products, paginated.
//	@tags			users
//	@security		ApiKeyAuth
//	@param			type		query	string	false	"Type of products, active/ended/expired"
//	@param			page		query	int		false	"Page Number"
//	@param			per_page	query	int		false	"Items per Page"
//	@produce		json
//	@success		200	{object}	users.GetProductsResponse	"Successfully retrieved"
//	@failure		400	{object}	shared.ErrorResponse		"When the request is invalid"
//	@failure		401	{object}	shared.ErrorResponse		"When the user is unauthenticated"
//	@failure		500	{object}	shared.ErrorResponse		"The server could not make the request"
//	@router			/users/me/products [get]
func (h *UsersHandler) GetMyProducts(g *gin.Context) {
	ctx := g.Request.Context()
	claimsAny, _ := g.Get("claims")
	claims := claimsAny.(*services.JWTSubject)
	query := GetMyProductsQuery{
		Page:    1,
		PerPage: 20,
		Type:    "active",
	}

	if err := g.ShouldBindQuery(&query); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "invalid query"})
		return
	}

	var state models.ProductState
	switch query.Type {
	case "ended":
		state = models.ProductStateEnded
	case "expired":
		state = models.ProductStateExpired
	default:
		state = models.ProductStateActive
	}

	products, err := h.ProductRepo.GetUserProducts(ctx, claims.UserID, state, query.PerPage, (query.Page-1)*query.PerPage)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "couldn't query for user products"})
		return
	}

	count, err := h.ProductRepo.CountUserProducts(ctx, claims.UserID, state)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "unable to count products"})
		return
	}

	response := GetProductsResponse{
		Data:       ranges.EachAddress(products, ToProductDTO),
		Total:      count,
		TotalPages: int(math.Ceil(float64(count) / float64(query.PerPage))),
		Page:       query.Page,
		PerPage:    query.PerPage,
	}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "query": query, "response": response})
	g.JSON(http.StatusOK, response)
}

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
	claimsAny, ok := g.Get("claims")
	if !ok {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": "no bearer authentication", "status": http.StatusUnauthorized})
		g.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse{Error: "no bearer authentication"})
		return
	}

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

// GetMyBids godoc
//
//	@summary		Gets my own bids.
//	@description	Retrieves information about your ongoing bids.
//	@tags			users
//	@produce		json
//	@security		ApiKeyAuth
//	@param			page			query		int							false	"Page Number"
//	@param			per_page		query		int							false	"Items per Page"
//	@param			includes_loss	query		boolean						false	"Includes non-winning auctions"
//	@success		200				{object}	users.GetProductsResponse	"Successful"
//	@failure		401				{object}	shared.ErrorResponse		"When unauthenticated"
//	@failure		500				{object}	shared.ErrorResponse		"The server could not complete the request"
//	@router			/users/me/bids [GET]
func (h *UsersHandler) GetMyBids(g *gin.Context) {
	claimsAny, _ := g.Get("claims")
	claims := claimsAny.(*services.JWTSubject)
	ctx := g.Request.Context()
	query := GetMyBidsQuery{
		Page:    1,
		PerPage: 20,
		Status:  "active",
	}

	if err := g.ShouldBindQuery(&query); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "invalid query"})
		return
	}

	products, err := h.ProductRepo.GetMyBids(ctx, claims.UserID, query.Status == "active", query.PerPage, (query.Page-1)*query.PerPage)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "couldn't query for user bids"})
		return
	}

	count, err := h.ProductRepo.CountMyBids(ctx, claims.UserID, query.Status == "active")
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "unable to count products"})
		return
	}

	response := GetProductsResponse{
		Data:       ranges.EachAddress(products, ToProductDTO),
		Total:      count,
		TotalPages: int(math.Ceil(float64(count) / float64(query.PerPage))),
		Page:       query.Page,
		PerPage:    query.PerPage,
	}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "query": query, "response": response})
	g.JSON(http.StatusOK, response)
}

// GetMyRatings godoc
//
//	@summary		Gets a list of ratings made by me.
//	@description	Gets all ratings and related products and feedbacks, made by me.
//	@tags			users
//	@security		ApiKeyAuth
//	@produce		json
//	@param			page		query		int							false	"Page number"
//	@param			per_page	query		int							false	"Items per page"
//	@success		200			{object}	users.GetRatingsResponse	"Success"
//	@failure		400			{object}	shared.ErrorResponse		"Invalid query"
//	@failure		401			{object}	shared.ErrorResponse		"Unauthenticated"
//	@failure		500			{object}	shared.ErrorResponse		"The server could not complete the request"
//	@router			/users/me/ratings [GET]
func (h *UsersHandler) GetMyRatings(g *gin.Context) {
	ctx := g.Request.Context()
	claimsAny, _ := g.Get("claims")
	claims := claimsAny.(*services.JWTSubject)
	query := GetMyProductsQuery{
		Page:    1,
		PerPage: 20,
	}

	if err := g.ShouldBindQuery(&query); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "invalid query"})
		return
	}

	ratings, err := h.RatingRepo.GetMyRatings(ctx, claims.UserID, query.PerPage, (query.Page-1)*query.PerPage)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "failed to read ratings"})
		return
	}

	count, err := h.RatingRepo.CountMyRatings(ctx, claims.UserID)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "unable to count products"})
		return
	}

	response := GetRatingsResponse{
		Data:       ranges.EachAddress(ratings, ToRatingDTO),
		Total:      count,
		TotalPages: int(math.Ceil(float64(count) / float64(query.PerPage))),
		Page:       query.Page,
		PerPage:    query.PerPage,
	}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "query": query, "response": response})
	g.JSON(http.StatusOK, response)
}

// GetMyRated godoc
//
//	@summary		Gets a list of ratings on me.
//	@description	Gets all ratings and related products and feedbacks, made on me.
//	@tags			users
//	@security		ApiKeyAuth
//	@produce		json
//	@param			page		query		int							false	"Page number"
//	@param			per_page	query		int							false	"Items per page"
//	@success		200			{object}	users.GetRatingsResponse	"Success"
//	@failure		400			{object}	shared.ErrorResponse		"Invalid query"
//	@failure		401			{object}	shared.ErrorResponse		"Unauthenticated"
//	@failure		500			{object}	shared.ErrorResponse		"The server could not complete the request"
//	@router			/users/me/rated [GET]
func (h *UsersHandler) GetMyRated(g *gin.Context) {
	ctx := g.Request.Context()
	claimsAny, _ := g.Get("claims")
	claims := claimsAny.(*services.JWTSubject)
	query := GetMyProductsQuery{
		Page:    1,
		PerPage: 20,
	}

	if err := g.ShouldBindQuery(&query); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "invalid query"})
		return
	}

	ratings, err := h.RatingRepo.GetMyReviewedRatings(ctx, claims.UserID, query.PerPage, (query.Page-1)*query.PerPage)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "failed to read ratings"})
		return
	}

	count, err := h.RatingRepo.CountMyReviewedRatings(ctx, claims.UserID)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "unable to count products"})
		return
	}

	response := GetRatingsResponse{
		Data:       ranges.EachAddress(ratings, ToRatingDTO),
		Total:      count,
		TotalPages: int(math.Ceil(float64(count) / float64(query.PerPage))),
		Page:       query.Page,
		PerPage:    query.PerPage,
	}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "query": query, "response": response})
	g.JSON(http.StatusOK, response)
}
