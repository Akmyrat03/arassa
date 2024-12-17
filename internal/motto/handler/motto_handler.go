package handler

import (
	"arassachylyk/internal/motto/model"
	"arassachylyk/internal/motto/service"
	handler "arassachylyk/pkg/response"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MottoHandler struct {
	service *service.MottoService
}

func NewYearHandler(service *service.MottoService) *MottoHandler {
	return &MottoHandler{service: service}
}

// AddMotto
// @Summary Add a new motto
// @Description Adds a new motto with translations and an uploaded image.
// @Tags Motto
// @Accept multipart/form-data
// @Produce json
// @Param name_tkm formData string true "Motto name in Turkmen"
// @Param name_eng formData string true "Motto name in English"
// @Param name_rus formData string true "Motto name in Russian"
// @Param image formData file true "Motto image file"
// @Success 200 {object} response.ErrorResponse "Successfully created motto"
// @Failure 400 {object} response.ErrorResponse "Bad request error message"
// @Failure 500 {object} response.ErrorResponse "Internal server error message"
// @Router /motto/add [post]
func (h *MottoHandler) AddMotto(c *gin.Context) {
	nameTkm := c.PostForm("name_tkm")
	nameEng := c.PostForm("name_eng")
	nameRus := c.PostForm("name_rus")

	translations := []model.Translation{
		{LangID: 1, Name: nameTkm},
		{LangID: 2, Name: nameEng},
		{LangID: 3, Name: nameRus},
	}

	image, err := c.FormFile("image")
	if err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	uploadDir := "./uploads/motto"

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}

	filepath := filepath.Join(uploadDir, image.Filename)

	if err := c.SaveUploadedFile(image, filepath); err != nil {
		handler.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	input := model.Motto{
		ImageURL:     filepath,
		Translations: translations,
	}

	id, err := h.service.Create(input)
	if err != nil {
		handler.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":      id,
		"message": "Successfully created motto",
	})
}

// DeleteMotto
// @Summary Delete a motto
// @Description Deletes a motto by ID and removes the associated image file if it exists.
// @Tags Motto
// @Accept json
// @Produce json
// @Param id path int true "Motto ID"
// @Success 200 {object} response.ErrorResponse "Successfully deleted"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 404 {object} response.ErrorResponse "Motto not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /motto/delete/{id} [delete]
func (h *MottoHandler) DeleteMotto(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	existingYear, err := h.service.GetByID(id)
	if err != nil {
		handler.NewErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		handler.NewErrorResponse(c, http.StatusInternalServerError, "Cannot delete year")
		return
	}

	if existingYear.ImageURL != "" {
		if err := os.Remove(existingYear.ImageURL); err != nil && !os.IsNotExist(err) {
			handler.NewErrorResponse(c, http.StatusInternalServerError, "Failed to delete image file")
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully deleted",
		"id":      id,
	})

}
