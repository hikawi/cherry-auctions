package categories

import (
	"time"

	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/pkg/slug"
)

type CategoryDTO struct {
	ID            uint           `json:"id"`
	Name          string         `json:"name"`
	Slug          string         `json:"slug"`
	ParentID      *uint          `json:"parent_id"`
	SubCategories []*CategoryDTO `json:"subcategories" swaggertype:"array,object"` // To make Swagger not dig too deep
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     *time.Time     `json:"deleted_at"`
}

// FromModel takes a Category model and translates it to a DTO.
// This function is NOT recursive.
func FromModel(m models.Category) CategoryDTO {
	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}

	return CategoryDTO{
		ID:            m.ID,
		Name:          m.Name,
		ParentID:      m.ParentID,
		Slug:          slug.Slugify(m.Name),
		SubCategories: make([]*CategoryDTO, 0),
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		DeletedAt:     deletedAt,
	}
}

type GetCategoriesResponse []CategoryDTO

// For creating a new category
type PostCategoryBody struct {
	Name     string `json:"name" binding:"min=2,max=200,required"`
	ParentID *uint  `json:"parent_id" binding:"gt=0,omitempty"`
}

// For returning the category created.
type PostCategoryResponse CategoryDTO

type PutCategoryBody struct {
	Name     *string `json:"name" binding:"min=2,max=200,required,omitnil"`
	ParentID *uint   `json:"parent_id" binding:"omitnil,gt=0"`
}
