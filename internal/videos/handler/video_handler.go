package handler

import (
	"arassachylyk/internal/videos/model"
	"arassachylyk/internal/videos/service"
	handler "arassachylyk/pkg/response"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	service *service.VideoService
}

func NewVideoHandler(service *service.VideoService) *VideoHandler {
	return &VideoHandler{service: service}
}

// UploadVideos
// @Summary      Upload videos with title translations
// @Description  Upload multiple videos with their title translations in different languages (e.g., Turkmen, English, Russian).
// @Tags         Videos
// @Accept       multipart/form-data
// @Produce      json
// @Param        title_tkm  formData  string  true   "Title in Turkmen"
// @Param        title_eng  formData  string  true   "Title in English"
// @Param        title_rus  formData  string  true   "Title in Russian"
// @Param        videos     formData  file    true   "Video files"
// @Success      200        {object} response.ErrorResponse "Successfully uploaded videos"
// @Failure      400        {object} response.ErrorResponse "Invalid form data or file size exceeds the limit"
// @Failure      500        {object} response.ErrorResponse "Failed to upload video"
// @Router       /videos/upload [post]
// @security     BearerAuth
func (h *VideoHandler) UploadVideos() gin.HandlerFunc {
	return func(c *gin.Context) {

		const maxUploadSize = 50 * 1024 * 1024 // 50MB limit for video files

		if c.Request.ContentLength > maxUploadSize {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds the 50MB limit"})
			return
		}

		// Extract translations from form data
		titleTkm := c.PostForm("title_tkm")
		titleEng := c.PostForm("title_eng")
		titleRus := c.PostForm("title_rus")

		translations := []model.Translation{
			{LangID: 1, Title: titleTkm},
			{LangID: 2, Title: titleEng},
			{LangID: 3, Title: titleRus},
		}

		// Process uploaded videos
		var videoPaths []string
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
			return
		}

		files := form.File["videos"]

		// Check file sizes and save files
		uploadDir := "./uploads/videos"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.MkdirAll(uploadDir, 0755)
		}

		for _, file := range files {
			if file.Size > maxUploadSize {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Video size exceeds the 50MB limit"})
				return
			}

			filePath := filepath.Join(uploadDir, file.Filename)
			if err := c.SaveUploadedFile(file, filePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload video"})
				return
			}

			videoPaths = append(videoPaths, filePath)
		}

		// Construct input model
		input := model.Title{
			Translations: translations,
			Videos:       videoPaths,
		}

		// Save to database
		id, err := h.service.UploadVideos(input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"id":      id,
			"message": "Successfully uploaded videos",
		})

	}
}

// DeleteVideos
// @Summary      Delete video and associated files
// @Description  Deletes a video by ID and removes associated files from the file system and database
// @Tags         Videos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Video Title ID"
// @Success      200  {object}  response.ErrorResponse "Videos deleted successfully"
// @Failure      400  {object}  response.ErrorResponse "Invalid ID"
// @Failure      500  {object}  response.ErrorResponse "Failed to delete files"
// @Router       /videos/delete/{id} [delete]
// @security BearerAuth
func (h *VideoHandler) DeleteVideos() gin.HandlerFunc {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
			return
		}

		// Step 1: Get file paths
		paths, err := h.service.GetVideoPaths(id)
		if err != nil {
			handler.NewErrorResponse(c, http.StatusNotFound, "Failed to retrieve videos")
			return
		}

		// Step 2: Delete files from the file system
		for _, path := range paths {
			err := os.Remove(path)
			if err != nil {
				handler.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
				continue
			}
		}

		// Step 3: Delete database record
		err = h.service.Delete(id)
		if err != nil {
			handler.NewErrorResponse(c, http.StatusInternalServerError, "Failed to delete videos")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Video deleted successfully",
		})

	}
}

// GetAllVideos
// @Summary Get videos by language
// @Description Retrieve a list of videos filtered by language using the lang_id query parameter
// @Tags Videos
// @Accept json
// @Produce json
// @Param lang_id query int true "Language ID (e.g., 1 for Turkmen, 2 for English, 3 for Russian)"
// @Success 200 {object} response.ErrorResponse "Successfully retrieved videos"
// @Failure 400 {object} response.ErrorResponse "Invalid lang_id"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /videos/all [get]
func (h *VideoHandler) GetAllVideos(c *gin.Context) {
	langID, err := strconv.Atoi(c.Query("lang_id"))
	if err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid lang_id")
		return
	}

	videos, err := h.service.GetAll(langID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"videos": videos})
}
