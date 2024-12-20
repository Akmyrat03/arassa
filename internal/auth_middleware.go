package internal

import (
	"arassachylyk/pkg/response"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.NewErrorResponse(c, http.StatusUnauthorized, "Token gereklidir")
			c.Abort()
			return
		}

		// Remove "Bearer " prefix from the token
		token = strings.TrimPrefix(token, "Bearer ")

		username, err := ValidateToken(token)
		if err != nil {
			response.NewErrorResponse(c, http.StatusUnauthorized, "Geçersiz token")
			c.Abort()
			return
		}

		// Add the username to the context
		c.Set("username", username)
		fmt.Println("Authorized by: ", username)

		c.Next()
	}
}
