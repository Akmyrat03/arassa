package middleware

import (
	"arassachylyk/internal/admin/model"
	"arassachylyk/internal/admin/repository"
	"arassachylyk/internal/admin/service"
	handler "arassachylyk/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AdminMiddleware struct {
	repo    *repository.AdminRepository
	service *service.AdminService
}

func NewAdminMiddleware(repo *repository.AdminRepository, service *service.AdminService) *AdminMiddleware {
	return &AdminMiddleware{
		repo:    repo,
		service: service,
	}
}

func (m *AdminMiddleware) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input model.Admin
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		if input.Username == "" || input.Password == "" {
			handler.NewErrorResponse(c, http.StatusBadRequest, "All fields are required")
			c.Abort()
			return
		}

		if len(input.Password) < 4 {
			handler.NewErrorResponse(c, http.StatusBadRequest, "Password must be at least 4 characters")
			c.Abort()
			return
		}

		// Check for existing user by username
		existingUser, err := m.repo.GetUserByField("username", input.Username)
		if err == nil && existingUser.Username != "" {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Username already exists",
			})
			c.Abort()
			return
		}

		// Create the user since it doesn't exist
		if _, err := m.service.CreateAdmin(&input); err != nil {
			handler.NewErrorResponse(c, http.StatusInternalServerError, "Could not create user")
			return
		}

		// Respond with a success message
		c.JSON(http.StatusCreated, gin.H{
			"message":  "User created successfully",
			"username": input.Username,
		})
		c.Next()
	}
}

// Login
// @Summary Login and get JWT token
// @Description Login by providing username and password as form data to get a JWT token for authentication.
// @Tags Login
// @Accept multipart/form-data
// @Produce json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {object} response.ErrorResponse "Login successful, returns admin username and JWT token"
// @Failure 400 {object} response.ErrorResponse "Missing or invalid input"
// @Failure 401 {object} response.ErrorResponse "Invalid username or password"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /admin/login [post]
func (m *AdminMiddleware) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		if username == "" || password == "" {
			handler.NewErrorResponse(c, http.StatusBadRequest, "All fields are required")
			return
		}

		token, err := m.service.GenerateToken(username, password)
		if err != nil {
			handler.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		admin, err := m.service.GetAdmin(username, service.GeneratePasswordHash(password))
		if err != nil {
			handler.NewErrorResponse(c, http.StatusUnauthorized, "username or password is incorrect")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"admin": admin.Username,
			"token": token,
		})
	}
}

func (m *AdminMiddleware) Profile() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			handler.NewErrorResponse(c, http.StatusUnauthorized, "Token gereklidir")
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		username, err := m.service.ValidateToken(token)
		if err != nil {
			handler.NewErrorResponse(c, http.StatusUnauthorized, "Geçersiz token")
			c.Abort()
			return
		}

		user, err := m.repo.GetUserByField("username", username)
		if err != nil {
			handler.NewErrorResponse(c, http.StatusNotFound, "Kullanıcı bulunamadı")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"username": user.Username,
		})
	}
}
