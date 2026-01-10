package shared

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type PaginationResponse struct {
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
}
