package chat

import (
	"mime/multipart"

	"luny.dev/cherryauctions/internal/routes/shared"
)

type ChatSessionResponse struct {
	Data       []shared.ChatSessionDTO `json:"data"`
	Total      int64                   `json:"total"`
	TotalPages int                     `json:"total_pages"`
	Page       int                     `json:"page"`
	PerPage    int                     `json:"per_page"`
}

type CreateChatSessionRequest struct {
	ProductID uint `json:"product_id"`
}

type PostMessageRequest struct {
	Content string                `form:"content" json:"content"`
	Image   *multipart.FileHeader `form:"image"`
}

type ChatMessageResponse struct {
	Data       []shared.ChatMessageDTO `json:"data"`
	Total      int64                   `json:"total"`
	TotalPages int                     `json:"total_pages"`
	Page       int                     `json:"page"`
	PerPage    int                     `json:"per_page"`
}
