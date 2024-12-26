package handler

import (
	"arassachylyk/internal/images/model"
	"arassachylyk/internal/images/service"
	"arassachylyk/pkg/consts"
	handler "arassachylyk/pkg/response"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ImageHandler struct {
	service *service.ImageService
}

func NewImageHandler(service *service.ImageService) *ImageHandler {
	return &ImageHandler{service: service}
}

// CreateImages
// @Summary Upload multiple images and create a title with translations
// @Description Upload multiple images with associated titles in different languages. Requires a valid JWT token
// @Tags Images
// @Accept multipart/form-data
// @Produce json
// @Param titleTurkmen formData string true "Title in Turkmen"
// @Param titleEnglish formData string true "Title in English"
// @Param titleRussian formData string true "Title in Russian"
// @Param images formData file true "Images to upload" multiple
// @Success 200 {object} response.ErrorResponse "Successfully added title and images"
// @Failure 400 {object} response.ErrorResponse "Invalid input or file size exceeds limit"
// @Failure 401 {object} response.ErrorResponse "Invalid or expired token"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /images [post]
// @security BearerAuth.
func (h *ImageHandler) CreateImages() gin.HandlerFunc {
	return func(c *gin.Context) {
		const maxUploadSize = 10 * 1024 * 1024 // 10MB

		if c.Request.ContentLength > maxUploadSize {
			handler.NewErrorResponse(c, http.StatusBadRequest, "10MB limit image size")
			return
		}

		titleTkm := c.PostForm("titleTurkmen")
		titleEng := c.PostForm("titleEnglish")
		titleRus := c.PostForm("titleRussian")

		translations := []model.Translation{
			{LangID: consts.LangIDTurkmen, Title: titleTkm},
			{LangID: consts.LangIDEnglish, Title: titleEng},
			{LangID: consts.LangIDRussian, Title: titleRus},
		}

		var images []string
		form, err := c.MultipartForm()
		if err != nil {
			handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid form data")
			return
		}

		files := form.File["images"]

		for _, file := range files {
			if file.Size > maxUploadSize {
				handler.NewErrorResponse(c, http.StatusBadRequest, "Image sizes is greater than 10MB")
				return
			}
		}

		uploadDir := "./uploads/images"

		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			if err := os.Mkdir(uploadDir, consts.DefaultPermissions); err != nil {
				fmt.Printf("Failed to create directory %s: %v\n", uploadDir, err)
				return
			}
		}

		for _, file := range files {
			filePath := filepath.Join(uploadDir, file.Filename)
			if err := c.SaveUploadedFile(file, filePath); err != nil {
				handler.NewErrorResponse(c, http.StatusInternalServerError, "Failed to upload image")
				return
			}

			images = append(images, filePath)

		}

		input := model.Title{
			Images:       images,
			Translations: translations,
		}

		id, err := h.service.Create(input)
		if err != nil {
			handler.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"id":      id,
			"message": "Successfully added title and images",
		})
	}
}

// DeleteImages
// @Summary Delete a title and its associated images
// @Description Delete a title and its associated images by ID. Requires a valid JWT token in the Authorization header.
// @Tags Images
// @Accept json
// @Produce json
// @Param id path int true "Title ID"
// @Success 200 {object} response.ErrorResponse "Successfully deleted title and its images"
// @Failure 400 {object} response.ErrorResponse "Invalid title ID"
// @Failure 401 {object} response.ErrorResponse "Invalid or expired token"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /images/{id} [delete]
// @security BearerAuth.
func (h *ImageHandler) DeleteImages() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			handler.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		imagePaths, err := h.service.GetImageByTitleID(id)
		if err != nil {
			handler.NewErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve image paths")
			return
		}

		for _, imagePath := range imagePaths {
			err := os.Remove(imagePath)
			if err != nil {
				handler.NewErrorResponse(c, http.StatusInternalServerError, "Cant delete images from uploads/images folder")
				continue
			}

		}
		err = h.service.Delete(id)
		if err != nil {
			handler.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully deleted title and its images",
		})

	}
}

// GetAllImages
// @Summary Get images by language
// @Description Retrieve a list of images filtered by language using the lang_id query parameter
// @Tags Images
// @Accept json
// @Produce json
// @Param lang_id query int true "Language ID (e.g., 1 for Turkmen, 2 for English, 3 for Russian)"
// @Success 200 {object} response.ErrorResponse "Successfully retrieved images"
// @Failure 400 {object} response.ErrorResponse "Invalid lang_id"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /images/all [get].
func (h *ImageHandler) GetAllImages(c *gin.Context) {
	langID, err := strconv.Atoi(c.Query("lang_id"))
	if err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid lang_id")
		return
	}

	images, err := h.service.GetAll(langID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"images": images})
}

func (h *ImageHandler) GetPaginatedImages(c *gin.Context) {
	langID, err := strconv.Atoi(c.Query("lang_id"))
	if err != nil || langID <= 0 {
		handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid or missing language ID")
		return
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid page number")
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid limit")
		return
	}

	images, err := h.service.GetPaginatedImg(langID, page, limit)
	if err != nil {
		handler.NewErrorResponse(c, http.StatusInternalServerError, "Could not get paginated images")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"images": images,
	})
}
