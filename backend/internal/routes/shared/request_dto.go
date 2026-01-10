package shared

type PaginationRequest struct {
	Page    int `form:"page" binding:"number,gt=0,omitempty" json:"page"`
	PerPage int `form:"per_page" binding:"number,gt=0,omitempty" json:"per_page"`
}
