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

type SSEvent struct {
	Type    string `json:"type"`
	Message any    `json:"message"`
}

// Hub manages all active SSE streams
type Hub struct {
	// Map of UserID -> Channel for messages
	Clients      map[uint]chan SSEvent
	SessionCache map[uint]*models.ChatSession
	mu           sync.Mutex
}

var hub = &Hub{
	Clients:      make(map[uint]chan SSEvent),
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
			ch <- SSEvent{
				Type:    "message",
				Message: shared.ToChatMessageDTO(chatMsg),
			}
		}
		h.safeSend(sess.BuyerID, SSEvent{
			Type:    "message",
			Message: shared.ToChatMessageDTO(chatMsg),
		})
		hub.mu.Unlock()
	}()
}

func (h *ChatHandler) SendTransactionChangeNotification(chatSessionID uint, transaction *models.Transaction) {
	go func() {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			hub.mu.Lock()
			sess, ok := hub.SessionCache[chatSessionID]

			if !ok {
				hub.mu.Unlock()
				newSession, err := h.chatSessionRepo.GetChatSessionByID(ctx, chatSessionID)
				if err != nil {
					return
				}
				hub.mu.Lock()
				hub.SessionCache[newSession.ID] = &newSession
				sess = &newSession
			}

			event := SSEvent{
				Type:    "transaction",
				Message: gin.H{"chat_session_id": chatSessionID, "transaction_status": transaction.TransactionStatus},
			}

			h.safeSend(sess.SellerID, event)
			h.safeSend(sess.BuyerID, event)

			hub.mu.Unlock()
		}()
	}()
}

// Helper function to prevent blocking
func (h *ChatHandler) safeSend(userID uint, event SSEvent) {
	if ch, ok := hub.Clients[userID]; ok {
		select {
		case ch <- event:
			// Success
		default:
			// Client channel is full or blocked; skip to keep system fast
			logging.LogRaw(logging.LOG_WARN, gin.H{"user_id": userID, "msg": "dropping event, client blocked"})
		}
	}
}
