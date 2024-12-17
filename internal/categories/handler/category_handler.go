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
// @Summary      Create a new category
// @Description  Create a new category with translations for Turkmen, English, and Russian
// @Tags         Categories
// @Accept       multipart/form-data
// @Produce      json
// @Param        category_tkm  formData  string  true  "Category name in Turkmen"
// @Param        category_eng  formData  string  true  "Category name in English"
// @Param        category_rus  formData  string  true  "Category name in Russian"
// @Success      200 {object} map[string]interface{} "Category created successfully"
// @Failure      400 {object} map[string]interface{} "Invalid input"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /categories/add [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
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

// GetAllCategoriesTKM
// @Summary Get all categories in Turkmen language
// @Description Retrieves all categories available in the Turkmen language.
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {object} response.ErrorResponse "List of categories"
// @Failure 400 {object} response.ErrorResponse "Bad request "
// @Router /categories/tkm [get]
func (h *CategoryHandler) GetAllCategoriesTKM(c *gin.Context) {
	categories, err := h.service.GetAllByLangID(1)
	if err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"kategoriyalar": categories,
	})
}

// GetAllCategoriesTKM
// @Summary Get all categories in English language
// @Description Retrieves all categories available in the English language.
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {object} response.ErrorResponse "List of categories"
// @Failure 400 {object} response.ErrorResponse "Bad request "
// @Router /categories/eng [get]
func (h *CategoryHandler) GetAllCategoriesENG(c *gin.Context) {
	categories, err := h.service.GetAllByLangID(2)
	if err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

// GetAllCategoriesTKM
// @Summary Get all categories in Russian language
// @Description Retrieves all categories available in the Russian language.
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {object} response.ErrorResponse "List of categories"
// @Failure 400 {object} response.ErrorResponse "Bad request "
// @Router /categories/rus [get]
func (h *CategoryHandler) GetAllCategoriesRUS(c *gin.Context) {
	categories, err := h.service.GetAllByLangID(3)
	if err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"категории": categories,
	})
}
