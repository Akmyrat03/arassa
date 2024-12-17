package handler

import (
	"arassachylyk/internal/images/model"
	"arassachylyk/internal/images/service"
	handler "arassachylyk/pkg/response"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type ImageHandler struct {
	service *service.ImageService
}

func NewImageHandler(service *service.ImageService) *ImageHandler {
	return &ImageHandler{service: service}
}

func (h *ImageHandler) CreateImages(c *gin.Context) {
	titleTkm := c.PostForm("title_tkm")
	titleEng := c.PostForm("title_eng")
	titleRus := c.PostForm("title_rus")

	translations := []model.Translation{
		{LangID: 1, Title: titleTkm},
		{LangID: 2, Title: titleEng},
		{LangID: 3, Title: titleRus},
	}

	var images []string
	form, err := c.MultipartForm()
	if err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid form data")
		return
	}

	files := form.File["images"]

	uploadDir := "./uploads/images"

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
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
