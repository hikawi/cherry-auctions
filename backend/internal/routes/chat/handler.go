package chat

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/logging"
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/routes/shared"
	"luny.dev/cherryauctions/internal/services"
	"luny.dev/cherryauctions/pkg/closer"
)

// GetChatSessions godoc
//
//	@summary		Gets all chat sessions.
//	@description	Gets all users chat sessions.
//	@param			page		query	int	false	"Page Number"
//	@param			per_page	query	int	false	"Items per Page"
//	@tags			chat
//	@produce		json
//	@security		ApiKeyAuth
//	@success		200	{object}	chat.ChatSessionResponse	"Successful query"
//	@failure		400	{object}	shared.ErrorResponse		"Bad request"
//	@failure		401	{object}	shared.ErrorResponse		"Unauthorized"
//	@failure		500	{object}	shared.ErrorResponse		"The server could not complete the request"
//	@router			/chat [GET]
func (h *ChatHandler) GetChatSessions(g *gin.Context) {
	ctx := g.Request.Context()
	claims, _ := g.Get("claims")
	sub := claims.(*services.JWTSubject)
	query := shared.PaginationRequest{
		Page:    1,
		PerPage: 20,
	}

	if err := g.ShouldBind(&query); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "invalid query"})
		return
	}

	sessions, err := h.chatSessionRepo.GetUserChatSessions(ctx, sub.UserID, query.PerPage, (query.Page-1)*query.PerPage)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "couldn't find user sessions"})
		return
	}

	count, err := h.chatSessionRepo.CountUserChatSessions(ctx, sub.UserID)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "couldn't count user sessions"})
		return
	}

	var dtos []shared.ChatSessionDTO
	for _, session := range sessions {
		dtos = append(dtos, shared.ToChatSessionDTO(&session))
	}

	response := ChatSessionResponse{
		Data:       dtos,
		Total:      count,
		TotalPages: int(math.Ceil(float64(count) / float64(query.PerPage))),
		Page:       query.Page,
		PerPage:    query.PerPage,
	}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "response": response})
	g.JSON(http.StatusOK, response)
}

// CreateChatSession godoc
//
//	@summary		Creates a new chat session for a product.
//	@description	Creates a new chat session for a product, only callable by the seller or the winner.
//	@tags			chat
//	@accept			json
//	@produce		json
//	@security		ApiKeyAuth
//	@param			body	body		chat.CreateChatSessionRequest	true	"Body"
//	@success		201		{object}	shared.MessageResponse			"Successfully created chat session"
//	@failure		400		{object}	shared.ErrorResponse			"Bad request"
//	@failure		401		{object}	shared.ErrorResponse			"Unauthorized"
//	@failure		403		{object}	shared.ErrorResponse			"No authority to create session for the product"
//	@failure		409		{object}	shared.ErrorResponse			"Failed to create chat session"
//	@failure		500		{object}	shared.ErrorResponse			"The server could not complete the request"
//	@router			/chat [POST]
func (h *ChatHandler) CreateChatSession(g *gin.Context) {
	ctx := g.Request.Context()
	claims, _ := g.Get("claims")
	sub := claims.(*services.JWTSubject)

	var body CreateChatSessionRequest
	if err := g.ShouldBind(&body); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "body": body})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	product, err := h.productRepo.GetProductByID(ctx, int(body.ProductID))
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "body": body})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "can't query for product"})
		return
	}

	if product.ProductState != models.ProductStateEnded || product.CurrentHighestBid == nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": "product doesn't end with a bid", "body": body})
		g.AbortWithStatusJSON(http.StatusForbidden, shared.ErrorResponse{Error: "product doesn't end with a bid"})
		return
	}

	if product.SellerID != sub.UserID && product.CurrentHighestBid.UserID != sub.UserID {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": "not a conversation member", "body": body})
		g.AbortWithStatusJSON(http.StatusForbidden, shared.ErrorResponse{Error: "not a conversation member"})
		return
	}

	if product.ChatSession != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusConflict, "error": "chat session already created", "body": body})
		g.AbortWithStatusJSON(http.StatusConflict, shared.ErrorResponse{Error: "already created session"})
		return
	}

	chatSession := models.ChatSession{
		ProductID: product.ID,
		SellerID:  product.SellerID,
		BuyerID:   product.CurrentHighestBid.UserID,
	}
	err = h.chatSessionRepo.CreateChatSession(ctx, &chatSession)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "body": body})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "failed to create session"})
		return
	}

	response := shared.MessageResponse{Message: "created session"}
	logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusOK, "body": body, "response": response})
	g.JSON(http.StatusOK, response)
}

// GetChatMessages godoc
//
//	@summary		Retrieves a list of chat messages in a channel
//	@description	Retrieves a list of chat messages in a channel for history reasons.
//	@tags			chat
//	@produce		json
//	@param			id			path	int	true	"Channel ID"
//	@param			page		query	int	false	"Page Number"
//	@param			per_page	query	int	false	"Items per Page"
//	@security		ApiKeyAuth
//	@failure		200	{object}	chat.ChatMessageResponse	"Successful"
//	@failure		400	{object}	shared.ErrorResponse		"Bad request"
//	@failure		401	{object}	shared.ErrorResponse		"Unauthorized"
//	@failure		403	{object}	shared.ErrorResponse		"Not a chat participant"
//	@failure		500	{object}	shared.ErrorResponse		"The server could not complete the request"
//	@router			/chat/{id} [GET]
func (h *ChatHandler) GetChatMessages(g *gin.Context) {
	ctx := g.Request.Context()
	claims, _ := g.Get("claims")
	sub := claims.(*services.JWTSubject)
	query := shared.PaginationRequest{
		Page:    1,
		PerPage: 50,
	}

	if err := g.ShouldBind(&query); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "invalid query"})
		return
	}

	id, err := strconv.ParseUint(g.Param("id"), 10, 0)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "invalid id"})
		return
	}

	session, err := h.chatSessionRepo.GetChatSessionByID(ctx, uint(id))
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "couldn't find chat session"})
		return
	}

	if sub.UserID != session.BuyerID && sub.UserID != session.SellerID {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": "not a participant of that chat", "query": query})
		g.AbortWithStatusJSON(http.StatusForbidden, shared.ErrorResponse{Error: "not a participant of that chat"})
		return
	}

	messages, err := h.chatSessionRepo.GetSessionChatMessages(ctx, session.ID, query.PerPage, (query.Page-1)*query.PerPage)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "couldn't find user messages"})
		return
	}

	count, err := h.chatSessionRepo.CountSessionChatMessages(ctx, session.ID)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "query": query})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "couldn't count user messages"})
		return
	}

	var dtos []shared.ChatMessageDTO
	for _, msg := range messages {
		dtos = append(dtos, shared.ToChatMessageDTO(&msg))
	}

	response := ChatMessageResponse{
		Data:       dtos,
		Total:      count,
		TotalPages: int(math.Ceil(float64(count) / float64(query.PerPage))),
		Page:       query.Page,
		PerPage:    query.PerPage,
	}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "response": response})
	g.JSON(http.StatusOK, response)
}

// PostChatMessage godoc
//
//	@summary		Posts a chat message to a channel.
//	@description	Posts a chat message to a channel the user may access.
//	@tags			chat
//	@accept			mpfd
//	@produce		json
//	@param			id	path		int						true	"Channel ID"
//	@success		201	{object}	shared.MessageResponse	"Successfully posted a message to a channel"
//	@failure		400	{object}	shared.ErrorResponse	"Bad channel ID, or invalid body"
//	@failure		401	{object}	shared.ErrorResponse	"Unauthorized"
//	@failure		403	{object}	shared.ErrorResponse	"Unknown channel ID"
//	@failure		413	{object}	shared.ErrorResponse	"Image too heavy"
//	@failure		500	{object}	shared.ErrorResponse	"The server failed to complete the request"
//	@router			/chat/{id} [post]
func (h *ChatHandler) PostChatMessage(g *gin.Context) {
	ctx := g.Request.Context()
	claims, _ := g.Get("claims")
	sub := claims.(*services.JWTSubject)

	id, err := strconv.ParseUint(g.Param("id"), 10, 0)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "invalid id"})
		return
	}

	session, err := h.chatSessionRepo.GetChatSessionByID(ctx, uint(id))
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "couldn't find chat session"})
		return
	}

	if sub.UserID != session.BuyerID && sub.UserID != session.SellerID {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": "not a participant of that chat"})
		g.AbortWithStatusJSON(http.StatusForbidden, shared.ErrorResponse{Error: "not a participant of that chat"})
		return
	}

	var body PostMessageRequest
	if err := g.ShouldBind(&body); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	var imageUrl *string
	if body.Image != nil {

		if body.Image.Size > (10 << 20) /* 10MB */ {
			logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": "image too large", "status": http.StatusRequestEntityTooLarge, "size": body.Image.Size})
			g.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, shared.ErrorResponse{Error: "max 10MB allowed"})
			return
		}

		file, _ := body.Image.Open()
		defer closer.CloseResources(file)

		// govips can load directly from an io.Reader, which is memory efficient
		img, err := vips.NewImageFromReader(file)
		if err != nil {
			logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
			g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "invalid image format"})
			return
		}
		defer img.Close()

		// Removed the part that crops and resizes
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
		key := fmt.Sprintf("images/%s.webp", hashString)
		err = h.s3Service.PutObject(ctx, key, bytes.NewReader(webpBuffer))
		if err != nil {
			logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusInternalServerError})
			g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "storage upload failed"})
			return
		}

		// Set it up on the chat message
		url := fmt.Sprintf("%s/%s", h.s3PermURL, key)
		imageUrl = &url
	}

	chatMsg := models.ChatMessage{
		SenderID:      sub.UserID,
		Content:       body.Content,
		ChatSessionID: session.ID,
		ImageURL:      imageUrl,
	}
	err = h.chatSessionRepo.CreateChatMessage(ctx, &chatMsg)

	// Populate before sending notification
	chatMsg.Sender = models.User{
		ID:    sub.UserID,
		Email: &sub.Email,
		Name:  &sub.Name,
	}

	h.SendNotification(&chatMsg)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "status": http.StatusInternalServerError})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "chat message creation failed"})
		return
	}
}

// GetChatStream godoc
//
// Courtesy of AI
//
//	@Summary		Opens a SSE stream to the chat channel
//	@Description	Establishes a persistent HTTP connection to receive real-time messages and receipts via Server-Sent Events.
//	@Tags			chat
//	@Produce		text/event-stream
//	@Security		ApiKeyAuth
//	@Param			token	query	string	true	"Authentication token"
//	@Router			/chat/stream [get]
func (h *ChatHandler) GetChatStream(g *gin.Context) {
	claims, _ := g.Get("claims")
	sub := claims.(*services.JWTSubject)

	// 1. Initialize the user's channel
	messageChan := make(chan SSEvent)

	hub.mu.Lock()
	hub.Clients[sub.UserID] = messageChan
	hub.mu.Unlock()

	// 2. Clean up when the user leaves
	defer func() {
		hub.mu.Lock()
		delete(hub.Clients, sub.UserID)
		hub.mu.Unlock()
		close(messageChan)
	}()

	// 3. Set SSE Headers and start the stream
	g.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-messageChan:
			if !ok {
				return false
			}
			// Send the JSON message as an SSE event
			g.SSEvent(msg.Type, msg.Message)
			return true
		case <-g.Request.Context().Done():
			// User disconnected
			return false
		case <-time.After(30 * time.Second):
			g.SSEvent("heartbeat", "ping")
			return true
		}
	})
}
