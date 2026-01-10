package chat

import (
	"context"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/logging"
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/routes/shared"
)

// Hub manages all active SSE streams
type Hub struct {
	// Map of UserID -> Channel for messages
	Clients      map[uint]chan shared.ChatMessageDTO
	SessionCache map[uint]*models.ChatSession
	mu           sync.Mutex
}

var hub = &Hub{
	Clients:      make(map[uint]chan shared.ChatMessageDTO),
	SessionCache: make(map[uint]*models.ChatSession),
}

func (h *ChatHandler) SendNotification(chatMsg *models.ChatMessage) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Read in the message to know where to route.
		sess, ok := hub.SessionCache[chatMsg.ChatSessionID]
		if !ok {
			// Fetch from db.
			newSession, err := h.chatSessionRepo.GetChatSessionByID(ctx, chatMsg.ChatSessionID)
			if err != nil {
				logging.LogRaw(logging.LOG_ERROR, gin.H{"error": err.Error()})
				return
			}

			hub.SessionCache[newSession.ID] = &newSession
			sess = &newSession
		}

		hub.mu.Lock()
		if ch, ok := hub.Clients[sess.SellerID]; ok {
			ch <- shared.ToChatMessageDTO(chatMsg)
		}
		if ch, ok := hub.Clients[sess.BuyerID]; ok {
			ch <- shared.ToChatMessageDTO(chatMsg)
		}
		defer hub.mu.Unlock()
	}()
}
