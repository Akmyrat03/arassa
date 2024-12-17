package routes

import (
	"arassachylyk/internal/contact/handler"
	"arassachylyk/internal/contact/repository"
	"arassachylyk/internal/contact/service"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func InitContactRoutes(router *gin.RouterGroup, db *sqlx.DB) {
	userRepo := repository.NewContactRepository(db)
	userService := service.NewContactService(userRepo)
	userHandler := handler.NewContactHandler(userService)

	userRoutes := router.Group("/contact")

	userRoutes.POST("/message", userHandler.SendMessage)

}
