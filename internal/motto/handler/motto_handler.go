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
// @Summary      Add a new motto
// @Description  Adds a new motto with translations in three languages (Turkmen, English, Russian) and an image upload. Requires a valid Bearer token.
// @Tags         Motto
// @Accept       multipart/form-data
// @Produce      json
// @Param        name_tkm  formData string true "Motto name in Turkmen"
// @Param        name_eng  formData string true "Motto name in English"
// @Param        name_rus  formData string true "Motto name in Russian"
// @Param        image     formData file   true "Motto image"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /motto/add [post]
// @security BearerAuth
func (h *MottoHandler) AddMotto() gin.HandlerFunc {
	return func(c *gin.Context) {

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
}

// DeleteMotto
// @Summary      Delete a motto
// @Description  Deletes a motto by its ID. Requires a valid Bearer token.
// @Tags         Motto
// @Accept       json
// @Produce      json
// @Param        id             path   int    true "Motto ID"
// @Success      200 {object} response.ErrorResponse "Successfully deleted"
// @Failure      400 {object} response.ErrorResponse "Invalid ID format"
// @Failure      404 {object} response.ErrorResponse "Cannot get motto by ID"
// @Failure      500 {object} response.ErrorResponse "Cannot delete year"
// @Router       /motto/delete/{id} [delete]
// @security BearerAuth
func (h *MottoHandler) DeleteMotto() gin.HandlerFunc {
	return func(c *gin.Context) {

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
}

// GetAllMottos
// @Summary      Get all mottos
// @Description  Retrieves all mottos with translations filtered by language ID. Requires a valid Bearer token.
// @Tags         Motto
// @Accept       json
// @Produce      json
// @Param        lang_id  query int true "Language ID"
// @Success      200 {object} []model.MottoResponse
// @Failure      400 {object} response.ErrorResponse "Invalid language ID format"
// @Failure      404 {object} response.ErrorResponse "No mottos found"
// @Failure      500 {object} response.ErrorResponse "Cannot retrieve mottos"
// @Router       /motto/all [get]
func (h *MottoHandler) GetAllMottos() gin.HandlerFunc {
	return func(c *gin.Context) {
		langID, err := strconv.Atoi(c.Query("lang_id"))
		if err != nil {
			handler.NewErrorResponse(c, http.StatusBadRequest, "Invalid language ID format")
			return
		}

		mottos, err := h.service.GetAllMottos(langID)
		if err != nil {
			handler.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		if len(mottos) == 0 {
			handler.NewErrorResponse(c, http.StatusNotFound, "No mottos found")
			return
		}

		c.JSON(http.StatusOK, mottos)
	}
}
