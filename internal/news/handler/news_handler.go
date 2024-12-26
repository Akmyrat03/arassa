package handler

import (
	"arassachylyk/internal"
	"arassachylyk/internal/news/model"
	"arassachylyk/internal/news/service"
	handler "arassachylyk/pkg/response"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type NewsHandler struct {
	service *service.NewsService
}

func NewHandler(service *service.NewsService) *NewsHandler {
	return &NewsHandler{service: service}
}

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
// @Router /news [post]
// @security BearerAuth.
func (h *NewsHandler) CreateNews() gin.HandlerFunc {
	return func(c *gin.Context) {

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
}

// GetAllNewsPagination
// @Summary Get all news in a specific language
// @Description Retrieves all news articles available in a specific language by language ID with pagination.
// @Tags News
// @Accept json
// @Produce json
// @Param lang_id query int true "Language ID"
// @Param page query int true "Page Number"
// @Param limit query int true "Limit"
// @Success 200 {object} response.ErrorResponse "List of news articles"
// @Failure 400 {object} response.ErrorResponse "Invalid Language ID or Pagination Parameters"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /news/all [get].
func (h *NewsHandler) GetAllNewsPagination(c *gin.Context) {
	langID, err := strconv.Atoi(c.Query("lang_id"))
	if err != nil || langID <= 0 {
		handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid Language ID")
		return
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid Page Number")
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid Limit")
		return
	}

	news, err := h.service.GetAllNewsByLangID(langID, page, limit)
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
// @Param        category_id query     int  true  "Category ID"
// @Param 		 lang_id query int true "Language ID"
// @Param		 page query int true "Page Number"
// @Param		 limit query int true "Limit"
// @Success      200 {object} map[string]interface{} "List of news"
// @Failure      400 {object} map[string]interface{} "Invalid input"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /news [get].
func (h *NewsHandler) GetAllNewsByLangAndCategory(c *gin.Context) {
	langID, err := strconv.Atoi(c.Query("lang_id"))
	if err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid language ID")
		return
	}

	categoryID, err := strconv.Atoi(c.Query("category_id"))
	if err != nil || categoryID <= 0 {
		handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	news, err := h.service.GetAllNewsByLangAndCategory(langID, categoryID, limit, page)
	if err != nil {
		handler.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"news": news,
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
// @Router /news/{id} [delete]
// @security BearerAuth.
func (h *NewsHandler) DeleteNews() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			handler.NewErrorResponse(c, http.StatusUnauthorized, "Token gereklidir")
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		username, err := internal.ValidateToken(token)
		if err != nil {
			handler.NewErrorResponse(c, http.StatusUnauthorized, "Geçersiz token")
			c.Abort()
			return
		}

		fmt.Println("Authorized by: ", username)

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
}
