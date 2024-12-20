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
// @Description Create a new category by providing its translations in multiple languages (requires a valid JWT token in the Authorization header)
// @Tags Categories
// @Accept multipart/form-data
// @Produce json
// @Param category_tkm formData string true "Category in Turkmen"
// @Param category_eng formData string true "Category in English"
// @Param category_rus formData string true "Category in Russian"
// @Success 200 {object} response.ErrorResponse "Category created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid input"
// @Failure 401 {object} response.ErrorResponse "Invalid or expired token"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /categories/add [post]
// @security BearerAuth
func (h *CategoryHandler) CreateCategory() gin.HandlerFunc {
	return func(c *gin.Context) {

		categoryTkm := c.PostForm("category_tkm")
		categoryEng := c.PostForm("category_eng")
		categoryRus := c.PostForm("category_rus")

		translations := []model.Translation{
			{Name: categoryTkm, LangID: 1},
			{Name: categoryEng, LangID: 2},
			{Name: categoryRus, LangID: 3},
		}

		req := model.CategoryReq{
			Translations: translations,
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
}

// DeleteCategory
// @Summary Delete a category
// @Description Delete a category by ID (requires valid JWT token in the Authorization header)
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} response.ErrorResponse "Successfully deleted category"
// @Failure 400 {object} response.ErrorResponse "Invalid category ID"
// @Failure 401 {object} response.ErrorResponse "Invalid or expired token"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /categories/delete/{id} [delete]
// @security BearerAuth
func (h *CategoryHandler) DeleteCategory() gin.HandlerFunc {
	return func(c *gin.Context) {

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
}

// GetAllCategories
// @Summary Get all categories by langID
// @Description Retrieves all categories by language ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id query int true "Language ID"
// @Success 200 {object} response.ErrorResponse "List of categories"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Router /categories/all [get]
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	categories, err := h.service.GetAllByLangID(id)
	if err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}
