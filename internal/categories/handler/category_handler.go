package handler

import (
	"arassachylyk/internal/categories/model"
	"arassachylyk/internal/categories/service"
	handler "arassachylyk/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// CreateCategory
// @Summary Create a new category
// @Description Create a new category with translations
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body model.CategoryReq true "Category data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /categories/add [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req model.Category
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.Create(req)
	if err != nil {
		handler.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":      id,
		"message": "Category created successfully",
	})

}

// UpdateCategory updates an existing category
// @Summary Update an existing category
// @Description Update the name of a category by ID
// @Tags Categories
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Category ID"
// @Param name formData string true "Updated category name"
// @Success 200 {object} response.ErrorResponse "Category updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid input data"
// @Failure 404 {object} response.ErrorResponse "Category not found"
// @Failure 500 {object} response.ErrorResponse "Could not update category"
// @Router /categories/update/{id} [put]
// func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid ID")
// 		return
// 	}

// 	name := c.PostForm("name")
// 	if name == "" {
// 		handler.NewErrorResponse(c, http.StatusBadRequest, "Category name is required")
// 		return
// 	}

// 	input := model.Category{
// 		ID:   id,
// 		Name: name,
// 	}

// 	err = h.service.Update(id, input)
// 	if err != nil {
// 		handler.NewErrorResponse(c, http.StatusInternalServerError, "Could not update category")
// 		return
// 	}

// 	c.JSON(http.StatusOK, map[string]interface{}{
// 		"message": "Category updated successfully",
// 	})
// }

// DeleteCategory deletes a category
// @Summary Delete a category
// @Description Deletes a category by its ID
// @Tags Categories
// @Param id path int true "Category ID"
// @Success 200 {object} response.ErrorResponse "Successfully deleted category"
// @Failure 400 {object} response.ErrorResponse "Invalid type id"
// @Failure 500 {object} response.ErrorResponse "Could not delete category"
// @Router /categories/delete/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid type id")
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		handler.NewErrorResponse(c, http.StatusInternalServerError, "Could not delete category")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully deleted category",
	})

}

// GetAllCategories retrieves all categories
// @Summary Get all categories
// @Description Retrieves a list of all categories
// @Tags Categories
// @Produce json
// @Success 200 {object} response.ErrorResponse "List of all categories"
// @Failure 500 {object} response.ErrorResponse "Failed to get all categories"
// // @Router /categories/view-all [get]
// func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
// 	categories, err := h.service.GetAll()
// 	if err != nil {
// 		handler.NewErrorResponse(c, http.StatusInternalServerError, "Failed to get all categories")
// 		return
// 	}

// 	c.JSON(http.StatusOK, map[string]interface{}{
// 		"categories": categories,
// 	})
// }
