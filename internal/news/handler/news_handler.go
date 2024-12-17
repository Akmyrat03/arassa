package handler

import (
	"arassachylyk/internal/news/model"
	"arassachylyk/internal/news/service"
	handler "arassachylyk/pkg/response"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NewsHandler struct {
	service *service.NewsService
}

func NewHandler(service *service.NewsService) *NewsHandler {
	return &NewsHandler{service: service}
}

// CreateNews handles the creation of news items with multilingual support.
// @Summary Create a news item
// @Description Creates a news item with category, image, and translations in Turkmen, English, and Russian
// @Tags News
// @Accept multipart/form-data
// @Produce json
// @Param category_id formData int true "Category ID"
// @Param image formData file true "Image file"
// @Param title_tkm formData string true "Title in Turkmen"
// @Param description_tkm formData string true "Description in Turkmen"
// @Param title_eng formData string true "Title in English"
// @Param description_eng formData string true "Description in English"
// @Param title_rus formData string true "Title in Russian"
// @Param description_rus formData string true "Description in Russian"
// @Success 200 {object} response.ErrorResponse "Successfully created news"
// @Failure 400 {object} response.ErrorResponse "Invalid input or bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /news/add-news [post]
func (h *NewsHandler) CreateNews(c *gin.Context) {
	var news model.News

	categoryID, err := strconv.Atoi(c.PostForm("category_id"))
	if err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	image, err := c.FormFile("image")
	if err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	uploadDir := "./uploads/news"

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}

	filepath := filepath.Join(uploadDir, image.Filename)

	if err := c.SaveUploadedFile(image, filepath); err != nil {
		handler.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	titleTKM := c.PostForm("title_tkm")
	descriptionTKM := c.PostForm("description_tkm")

	titleENG := c.PostForm("title_eng")
	descriptionENG := c.PostForm("description_eng")

	titleRUS := c.PostForm("title_rus")
	descriptionRUS := c.PostForm("description_rus")

	translations := []model.Translation{
		{LangID: 1, Title: titleTKM, Description: descriptionTKM},
		{LangID: 2, Title: titleENG, Description: descriptionENG},
		{LangID: 3, Title: titleRUS, Description: descriptionRUS},
	}

	news = model.News{
		CategoryID:   categoryID,
		ImageURL:     filepath,
		Translations: translations,
	}

	id, err := h.service.Create(news)
	if err != nil {
		handler.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":      id,
		"message": "News created successfully",
		"image":   filepath,
	})
}

// DeleteNews deletes a news
// @Summary Delete a news
// @Description Deletes a news by its ID
// @Tags News
// @Param id path int true "News ID"
// @Success 200 {object} response.ErrorResponse "News deleted successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Could not delete news"
// @Router /news/delete/{id} [delete]
func (h *NewsHandler) DeleteNews(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	existingNews, err := h.service.GetNewsByID(id)
	if err != nil {
		handler.NewErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	if err := h.service.Delete(id); err != nil {
		handler.NewErrorResponse(c, http.StatusInternalServerError, "Could not delete news")
		return
	}

	if existingNews.ImageURL != "" {
		if err := os.Remove(existingNews.ImageURL); err != nil && !os.IsNotExist(err) {
			handler.NewErrorResponse(c, http.StatusInternalServerError, "Failed to delete image file")
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "News deleted successfully",
		"id":      id,
	})
}

// GetAllNewsTKM
// @Summary      Get all news in Turkmen language
// @Description  Retrieves a list of all news with titles, descriptions, categories, and images in Turkmen language.
// @Tags         News
// @Accept       json
// @Produce      json
// @Success      200 {object} response.ErrorResponse "Successfully get all news in Turkmen language"
// @Failure      500 {object} response.ErrorResponse "Internal server error"
// @Router       /news/tkm [get]
func (h *NewsHandler) GetAllNewsTKM(c *gin.Context) {
	news, err := h.service.GetAllNewsByLangID(1)
	if err != nil {
		handler.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"news": news,
	})
}

// GetAllNewsENG
// @Summary      Get all news in English language
// @Description  Retrieves a list of all news with titles, descriptions, categories, and images in Turkmen language.
// @Tags         News
// @Accept       json
// @Produce      json
// @Success      200 {object} response.ErrorResponse "Successfully get all news in English language"
// @Failure      500 {object} response.ErrorResponse "Internal server error"
// @Router       /news/eng [get]
func (h *NewsHandler) GetAllNewsENG(c *gin.Context) {
	news, err := h.service.GetAllNewsByLangID(2)
	if err != nil {
		handler.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"news": news,
	})
}

// GetAllNewsENG
// @Summary      Get all news in Russian language
// @Description  Retrieves a list of all news with titles, descriptions, categories, and images in Turkmen language.
// @Tags         News
// @Accept       json
// @Produce      json
// @Success      200 {object} response.ErrorResponse "Successfully get all news in Russian language"
// @Failure      500 {object} response.ErrorResponse "Internal server error"
// @Router       /news/rus [get]
func (h *NewsHandler) GetAllNewsRUS(c *gin.Context) {
	news, err := h.service.GetAllNewsByLangID(3)
	if err != nil {
		handler.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"news": news,
	})
}

// GetAllNewsByLangAndCategory
// @Summary      Get all news by language and category
// @Description  Fetch all news based on language ID and category ID
// @Tags         News
// @Accept       json
// @Produce      json
// @Param        lang_id     query     int  true  "Language ID (e.g., 1 for Turkmen, 2 for Russian)"
// @Param        category_id query     int  true  "Category ID"
// @Success      200 {object} map[string]interface{} "List of news"
// @Failure      400 {object} map[string]interface{} "Invalid input"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /news/category [get]
func (h *NewsHandler) GetAllNewsByLangAndCategory(c *gin.Context) {
	langID, err := strconv.Atoi(c.Query("lang_id"))
	if err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid language ID")
		return
	}

	categoryID, err := strconv.Atoi(c.Query("category_id"))
	if err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	news, err := h.service.GetAllNewsByLangAndCategory(langID, categoryID)
	if err != nil {
		handler.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"news": news,
	})
}
