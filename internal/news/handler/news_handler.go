package handler

import (
	"arassachylyk/internal/news/model"
	"arassachylyk/internal/news/service"
	handler "arassachylyk/pkg/response"
	"encoding/json"
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

	// var translations []model.Translation
	// if err := c.ShouldBindJSON(&translations); err != nil {
	// 	handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid translations")
	// 	return
	// }

	translationsStr := c.PostForm("translations")
	var translations []model.Translation
	if err := json.Unmarshal([]byte(translationsStr), &translations); err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid translations")
		return
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

// GetAllNews
// @Summary Get all news
// @Description Get a list of all news
// @Tags News
// @Accept json
// @Produce json
// @Success 200 {object} response.ErrorResponse "Success"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Router /news/all [get]
func (h *NewsHandler) GetAllNews(c *gin.Context) {
	news, err := h.service.GetAll()
	if err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"news": news,
	})
}

func (h *NewsHandler) GetNewsByCategoryID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	news, err := h.service.GetByCategoryID(id)
	if err != nil {
		handler.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"news": news,
	})
}
