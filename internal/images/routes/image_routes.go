package routes

import (
	"arassachylyk/internal"
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
	imgRoutes.Use(internal.AuthMiddleware())
	{
		imgRoutes.POST("/", imgHandler.CreateImages())
		imgRoutes.DELETE("/:id", imgHandler.DeleteImages())
	}
	router.GET("/images/all", imgHandler.GetAllImages)
	router.GET("/images", imgHandler.GetPaginatedImages)

}
