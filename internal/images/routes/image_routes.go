package routes

import (
	"arassachylyk/internal/images/handler"
	"arassachylyk/internal/images/repository"
	"arassachylyk/internal/images/service"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func InitImageRoutes(router *gin.RouterGroup, db *sqlx.DB) {
	imgRepo := repository.NewImageRepository(db)
	imgService := service.NewImageService(imgRepo)
	imgHandler := handler.NewImageHandler(imgService)

	imgRoutes := router.Group("/images")

	imgRoutes.POST("/add", imgHandler.CreateImages)
}
