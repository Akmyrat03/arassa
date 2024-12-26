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

	router.POST("/contact", userHandler.SendMessage)

}
