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

// AddYear adds a new year
// @Summary Add a new year
// @Description Add a new year entry with an image
// @Tags Motto
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Motto"
// @Param language_id formData int true "Language ID"
// @Param image formData file true "Image file"
// @Success 200 {object} response.ErrorResponse "Successfully created motto"
// @Failure 400 {object} response.ErrorResponse "Invalid input"
// @Failure 500 {object} response.ErrorResponse "Could not create motto"
// @Router /motto/add [post]
func (h *MottoHandler) AddYear(c *gin.Context) {
	name := c.PostForm("name")
	languageId, err := strconv.Atoi(c.PostForm("language_id"))
	if err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	image, err := c.FormFile("image")
	if err != nil {
		handler.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	uploadDir := "./uploads/years"

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}

	filepath := filepath.Join(uploadDir, image.Filename)

	if err := c.SaveUploadedFile(image, filepath); err != nil {
		handler.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	input := model.Motto{
		Name:       name,
		LanguageID: languageId,
		ImageURL:   filepath,
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

// DeleteYear deletes a year
// @Summary Delete a year entry
// @Description Delete a year entry by its ID
// @Tags Motto
// @Produce json
// @Param id path int true "Motto ID"
// @Success 200 {object} response.ErrorResponse "Successfully deleted"
// @Failure 400 {object} response.ErrorResponse "Invalid ID"
// @Failure 404 {object} response.ErrorResponse "Year not found"
// @Failure 500 {object} response.ErrorResponse "Failed to delete year or image file"
// @Router /motto/delete/{id} [delete]
func (h *MottoHandler) DeleteYear(c *gin.Context) {
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
