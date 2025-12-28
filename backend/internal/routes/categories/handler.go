package categories

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/logging"
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/routes/shared"
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
	ctx := g.Request.Context()

	categories, err := h.CategoryRepo.GetActiveCategories(ctx)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "can't fetch categories"})
		return
	}

	response := buildTree(categories)
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "response": response})
	g.JSON(http.StatusOK, response)
}

// PostCategories POST /categories
//
//	@summary	Creates a category
//	@tags		categories
//	@accept		json
//	@produce	json
//	@security	ApiKeyAuth
//	@param		data	body		categories.PostCategoryBody	true	"New category data"
//	@success	201		{object}	categories.CategoryDTO		"Category was successfully created"
//	@failure	400		{object}	shared.ErrorResponse		"Category body data was invalid"
//	@Failure	401		{object}	shared.ErrorResponse		"Unauthorized"
//	@Failure	403		{object}	shared.ErrorResponse		"Not enough permissions"
//	@failure	500		{object}	shared.ErrorResponse		"The server could not complete the request"
//	@router		/categories [POST]
func (h *CategoriesHandler) PostCategories(g *gin.Context) {
	ctx := g.Request.Context()

	var body PostCategoryBody
	if err := g.ShouldBind(&body); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "body": body})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	category := models.Category{
		Name:     body.Name,
		ParentID: body.ParentID,
	}
	if err := h.CategoryRepo.SaveCategory(ctx, &category); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "body": body})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "internal server error"})
		return
	}

	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "body": body, "response": FromModel(category)})
	g.JSON(http.StatusCreated, FromModel(category))
}

// PutCategories godoc
//
//	@Summary		Update a category record
//	@Description	Update a category records by updating it partially
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@security		ApiKeyAuth
//	@Param			id		path		int							true	"Category ID"	Format(int64)
//	@param			data	body		categories.PutCategoryBody	true	"New updated category data"
//	@Success		200		{object}	categories.CategoryDTO		"Category was successfully changed, returns old category"
//	@Failure		400		{object}	shared.ErrorResponse		"Category body data was invalid"
//	@Failure		401		{object}	shared.ErrorResponse		"Unauthorized"
//	@Failure		403		{object}	shared.ErrorResponse		"Not enough permissions"
//	@Failure		404		{object}	shared.ErrorResponse		"Category could not be found"
//	@Failure		500		{object}	shared.ErrorResponse		"The server could not complete the request"
//	@Router			/categories/{id} [put]
func (h *CategoriesHandler) PutCategories(g *gin.Context) {
	ctx := g.Request.Context()

	paramId := g.Param("id")
	id, err := strconv.ParseInt(paramId, 10, 0)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "id": paramId})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	var body PutCategoryBody
	if err := g.ShouldBind(&body); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "body": body})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	category, err := h.CategoryRepo.GetCategoryByID(ctx, uint(id))
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusNotFound, "error": err.Error(), "body": body, "id": paramId})
		g.AbortWithStatusJSON(http.StatusNotFound, shared.ErrorResponse{Error: "category not found"})
		return
	}

	_, err = h.CategoryRepo.UpdateCategory(ctx, uint(id), body.Name, body.ParentID)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "body": body, "id": paramId})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "couldn't save category"})
		return
	}

	res := FromModel(category)
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "response": res, "body": body, "id": paramId})
	g.JSON(http.StatusOK, res)
}

// DeleteCategories godoc
//
//	@Summary		Soft-delete a category record
//	@Description	Soft-deletes a category record by setting a flag.
//	@Tags			categories
//	@Produce		json
//	@security		ApiKeyAuth
//	@Param			id	path	int	true	"Category ID"	Format(int64)
//	@Success		204	"Category was deleted"
//	@Failure		400	{object}	shared.ErrorResponse	"Category body data was invalid"
//	@Failure		401	{object}	shared.ErrorResponse	"Unauthorized"
//	@Failure		403	{object}	shared.ErrorResponse	"Not enough permissions"
//	@Failure		404	{object}	shared.ErrorResponse	"Category could not be found"
//	@Failure		500	{object}	shared.ErrorResponse	"The server could not complete the request"
//	@Router			/categories/{id} [delete]
func (h *CategoriesHandler) DeleteCategories(g *gin.Context) {
	ctx := g.Request.Context()

	paramId := g.Param("id")
	id, err := strconv.ParseInt(paramId, 10, 0)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "id": paramId})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	_, err = h.CategoryRepo.SoftDeleteCategory(ctx, uint(id))
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "id": paramId})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "couldn't delete category"})
		return
	}

	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusNoContent, "id": paramId})
	g.Status(http.StatusNoContent)
}
