package categories

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/models"
	"luny.dev/cherryauctions/repositories"
	"luny.dev/cherryauctions/routes/shared"
	"luny.dev/cherryauctions/utils"
)

// buildTree builds a recursively nested tree before sending back to the frontend.
func buildTree(categories []models.Category) any {
	nodes := make(map[uint]*CategoryDTO)
	children := make(map[uint]uint)

	// First pass to convert to nodes and pointers.
	for _, cat := range categories {
		if cat.ParentID != nil {
			children[cat.ID] = *cat.ParentID
		}

		dto := FromModel(cat)
		nodes[cat.ID] = &dto
	}

	// Now we bind the tree's children.
	for child, parent := range children {
		nodes[parent].SubCategories = append(nodes[parent].SubCategories, nodes[child])
	}

	// We build the list now.
	response := make(GetCategoriesResponse, 0)
	for _, cat := range categories {
		if cat.ParentID == nil {
			response = append(response, *nodes[cat.ID])
		}
	}

	return response
}

// GetCategories GET /categories
//
//	@summary		Retrieves a recursively nested list of categories.
//	@description	This endpoint retrieves all categories in a fully recursive nested list, for easily displaying on the frontend.
//	@tags			categories
//	@produce		json
//	@success		200	{object}	categories.GetCategoriesResponse	"List was successfully fetched"
//	@failure		500	{object}	shared.ErrorResponse				"The request could not be completed due to server faults"
//	@router			/categories [GET]
func (h *CategoriesHandler) GetCategories(g *gin.Context) {
	ctx := context.Background()
	categoriesRepo := repositories.CategoryRepository{DB: h.DB, Context: ctx}

	categories, err := categoriesRepo.GetActiveCategories()
	if err != nil {
		utils.LogMessage(g, utils.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "can't fetch categories"})
	}

	response := buildTree(categories)

	utils.LogMessage(g, utils.LOG_INFO, gin.H{"status": http.StatusOK, "response": response})
	g.JSON(http.StatusOK, response)
}
