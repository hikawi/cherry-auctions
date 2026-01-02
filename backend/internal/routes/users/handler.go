package users

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"net/http"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/logging"
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/routes/shared"
	"luny.dev/cherryauctions/internal/services"
	"luny.dev/cherryauctions/pkg/closer"
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

// PostAvatar godoc
//
//	@summary		Changes my avatar.
//	@description	Uploads and changes my avatar.
//	@tags			users
//	@produce		json
//	@security		ApiKeyAuth
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			avatar	formData	file						true	"account image"
//	@success		200		{object}	users.PostAvatarResponse	"When success"
//	@failure		400		{object}	shared.ErrorResponse		"Invalid request"
//	@failure		401		{object}	shared.ErrorResponse		"When unauthenticated"
//	@failure		413		{object}	shared.ErrorResponse		"When the image is too big"
//	@failure		500		{object}	shared.ErrorResponse		"The request could not be completed due to server faults"
//	@router			/users/avatar [POST]
func (h *UsersHandler) PostAvatar(g *gin.Context) {
	ctx := g.Request.Context()
	claimsAny, _ := g.Get("claims")
	claims := claimsAny.(*services.JWTSubject)

	var body PostAvatarRequest
	if err := g.ShouldBind(&body); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	if body.Avatar.Size > (10 << 20) /* 10MB */ {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": "image too large", "status": http.StatusRequestEntityTooLarge, "size": body.Avatar.Size})
		g.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, shared.ErrorResponse{Error: "max 10MB allowed"})
		return
	}

	file, _ := body.Avatar.Open()
	defer closer.CloseResources(file)

	// govips can load directly from an io.Reader, which is memory efficient
	img, err := vips.NewImageFromReader(file)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "invalid image format"})
		return
	}
	defer img.Close()

	// .Thumbnail is the most optimized way to resize in vips
	// InterestingAttention uses an algorithm to find the most "interesting" part (usually a face)
	err = img.Thumbnail(256, 256, vips.InterestingAttention)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusInternalServerError})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "processing failed"})
		return
	}

	ep := vips.NewWebpExportParams()
	ep.Quality = 75
	ep.StripMetadata = true

	webpBuffer, _, err := img.ExportWebp(ep)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusInternalServerError})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "encoding failed"})
		return
	}

	// 6. Upload to S3
	hash := sha256.Sum256(webpBuffer)
	hashString := hex.EncodeToString(hash[:])
	key := fmt.Sprintf("avatars/%s.webp", hashString)
	err = h.S3Service.PutObject(ctx, key, bytes.NewReader(webpBuffer))
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusInternalServerError})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "storage upload failed"})
		return
	}

	// Set it up on the user.
	avatarURL := fmt.Sprintf("%s/%s", h.S3PermURL, key)
	_, err = h.UserRepo.UpdateAvatarURL(ctx, claims.UserID, avatarURL)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusInternalServerError})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "can't update avatar url"})
		return
	}

	response := PostAvatarResponse{
		AvatarURL: avatarURL,
	}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "response": response})
	g.JSON(http.StatusOK, response)
}

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
//	@router			/users/profile [PUT]
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
