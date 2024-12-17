package handler

import (
	"arassachylyk/internal/contact/model"
	"arassachylyk/internal/contact/service"
	handler "arassachylyk/pkg/response"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContactHandler struct {
	Service *service.ContactService
}

func NewContactHandler(service *service.ContactService) *ContactHandler {
	return &ContactHandler{Service: service}
}

// SendMessage
// @Summary Send a contact message
// @Description Send a message from a contact form
// @Tags Contact
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Param name formData string true "Name of the sender"
// @Param email formData string true "Email of the sender"
// @Param message formData string true "Message content"
// @Param phone_number formData string true "Phone number of the sender"
// @Success 200 {object} response.ErrorResponse "Mesaj başarıyla gönderildi"
// @Failure 400 {object} response.ErrorResponse "All fields are required"
// @Failure 500 {object} response.ErrorResponse "Mesaj gönderilemedi"
// @Router /contact/message [post]
func (h *ContactHandler) SendMessage(c *gin.Context) {
	var message model.ContactMessage

	name := c.PostForm("name")
	email := c.PostForm("email")
	messageContent := c.PostForm("message")
	phone_number := c.PostForm("phone_number")

	if name == "" || email == "" || phone_number == "" || messageContent == "" {
		handler.NewErrorResponse(c, http.StatusBadRequest, "All fields are required")
		return
	}

	message = model.ContactMessage{
		Name:        name,
		Email:       email,
		Message:     messageContent,
		PhoneNumber: phone_number,
	}

	ctx := context.Background()
	if err := h.Service.SendMessage(ctx, message); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mesaj başarıyla gönderildi"})
}
