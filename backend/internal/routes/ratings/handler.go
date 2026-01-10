package ratings

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/logging"
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/routes/shared"
	"luny.dev/cherryauctions/internal/services"
)

func (h *RatingHandler) Claims(g *gin.Context) *services.JWTSubject {
	claimsAny, _ := g.Get("claims")
	return claimsAny.(*services.JWTSubject)
}

// PostRating godoc
//
//	@summary		Creates a rating.
//	@description	Posts a rating and feedback onto a users profile.
//	@tags			ratings
//	@produce		json
//	@accept			json
//	@param			rating	body		ratings.PostRatingBody	true	"Rating context"
//	@success		201		{object}	shared.MessageResponse	"Successfully rated person"
//	@failure		400		{object}	shared.ErrorResponse	"Invalid reviewee ID or invalid format"
//	@failure		401		{object}	shared.ErrorResponse	"Unauthorized"
//	@failure		500		{object}	shared.ErrorResponse	"The server failed to complete the request"
//	@router			/ratings [POST]
func (h *RatingHandler) PostRating(g *gin.Context) {
	ctx := g.Request.Context()
	claims := h.Claims(g)
	var body PostRatingBody
	if err := g.ShouldBind(&body); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	if body.RevieweeID == claims.UserID {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": "can't rate self"})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "can't rate self"})
		return
	}

	rating := models.Rating{
		ReviewerID: claims.UserID,
		RevieweeID: body.RevieweeID,
		Rating:     body.Rating,
		Feedback:   body.Feedback,
	}
	err := h.ratingRepo.CreateRating(ctx, &rating)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "can't create rating"})
		return
	}

	response := shared.MessageResponse{Message: "created rating"}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusCreated, "body": body, "response": response})
	g.JSON(http.StatusCreated, response)
}

// PutRating godoc
//
//	@summary		Edits a rating.
//	@description	Edits a rating and feedback onto a users profile.
//	@tags			ratings
//	@produce		json
//	@accept			json
//	@param			id		path		int						true	"Rating ID"
//	@param			rating	body		ratings.PutRatingBody	true	"Rating context"
//	@success		200		{object}	shared.MessageResponse	"Successfully edited rating"
//	@failure		400		{object}	shared.ErrorResponse	"Invalid reviewee ID or invalid format"
//	@failure		401		{object}	shared.ErrorResponse	"Unauthorized"
//	@failure		403		{object}	shared.ErrorResponse	"Not your rating"
//	@failure		404		{object}	shared.ErrorResponse	"Rating ID not found"
//	@failure		500		{object}	shared.ErrorResponse	"The server failed to complete the request"
//	@router			/ratings/{id} [PUT]
func (h *RatingHandler) PutRating(g *gin.Context) {
	ctx := g.Request.Context()
	claims := h.Claims(g)
	var body PutRatingBody
	if err := g.ShouldBind(&body); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	paramId := g.Param("id")
	id, err := strconv.ParseUint(paramId, 10, 0)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	if id == 0 {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": "invalid id"})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	rating, err := h.ratingRepo.GetRatingByID(ctx, uint(id))
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": "invalid id"})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	if rating.RevieweeID == claims.UserID {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": "can't rate self"})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "can't rate self"})
		return
	}

	if rating.ReviewerID != claims.UserID {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": "not your review"})
		g.AbortWithStatusJSON(http.StatusForbidden, shared.ErrorResponse{Error: "not your review"})
		return
	}

	err = h.ratingRepo.UpdateRating(ctx, rating.ID, body.Rating, body.Feedback)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "can't edit rating"})
		return
	}

	response := shared.MessageResponse{Message: "edited rating"}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "body": body, "response": response})
	g.JSON(http.StatusOK, response)
}

// DeleteRating godoc
//
//	@summary		Deletes a rating.
//	@description	Deletes a rating.
//	@tags			ratings
//	@produce		json
//	@param			id	path		int						true	"Rating ID"
//	@success		200	{object}	shared.MessageResponse	"Successfully deleted rating"
//	@failure		400	{object}	shared.ErrorResponse	"Invalid reviewee ID or invalid format"
//	@failure		401	{object}	shared.ErrorResponse	"Unauthorized"
//	@failure		403	{object}	shared.ErrorResponse	"Not your rating"
//	@failure		404	{object}	shared.ErrorResponse	"Rating ID not found"
//	@failure		500	{object}	shared.ErrorResponse	"The server failed to complete the request"
//	@router			/ratings/{id} [DELETE]
func (h *RatingHandler) DeleteRating(g *gin.Context) {
	ctx := g.Request.Context()
	claims := h.Claims(g)

	paramId := g.Param("id")
	id, err := strconv.ParseUint(paramId, 10, 0)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	if id == 0 {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": "invalid id"})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	rating, err := h.ratingRepo.GetRatingByID(ctx, uint(id))
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": "invalid id"})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	if rating.RevieweeID == claims.UserID {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": "can't rate self"})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "can't rate self"})
		return
	}

	if rating.ReviewerID != claims.UserID {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": "not your review"})
		g.AbortWithStatusJSON(http.StatusForbidden, shared.ErrorResponse{Error: "not your review"})
		return
	}

	err = h.ratingRepo.DeleteRatingByID(ctx, rating.ID)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "can't delete rating"})
		return
	}

	response := shared.MessageResponse{Message: "deleted rating"}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "response": response})
	g.JSON(http.StatusOK, response)
}
