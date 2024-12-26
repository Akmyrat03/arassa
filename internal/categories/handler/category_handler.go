package handler

import (
	"arassachylyk/internal/categories/model"
	"arassachylyk/internal/categories/service"
	"arassachylyk/pkg/consts"
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
// @Param categoryTurkmen formData string true "Category in Turkmen"
// @Param categoryEnglish formData string true "Category in English"
// @Param categoryRussian formData string true "Category in Russian"
// @Success 200 {object} response.ErrorResponse "Category created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid input"
// @Failure 401 {object} response.ErrorResponse "Invalid or expired token"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /categories [post]
// @security BearerAuth.
func (h *CategoryHandler) CreateCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryTkm := c.PostForm("categoryTurkmen")
		categoryEng := c.PostForm("categoryEnglish")
		categoryRus := c.PostForm("categoryRussian")

		translations := []model.Translation{
			{Name: categoryTkm, LangID: consts.LangIDTurkmen},
			{Name: categoryEng, LangID: consts.LangIDEnglish},
			{Name: categoryRus, LangID: consts.LangIDRussian},
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
// @Router /categories/{id} [delete]
// @security BearerAuth.
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
// @Router /categories [get].
func (h *CategoryHandler) GetAllCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		langID, err := strconv.Atoi(c.Query("lang_id"))
		if err != nil || langID <= 0 {
			handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid ID")
			return
		}

		categories, err := h.service.GetAllByLangID(langID)
		if err != nil {
			handler.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"categories": categories,
		})
	}
}
