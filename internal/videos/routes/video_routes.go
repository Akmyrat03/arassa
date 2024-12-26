package routes

import (
	"arassachylyk/internal"
	"arassachylyk/internal/videos/handler"
	"arassachylyk/internal/videos/repository"
	"arassachylyk/internal/videos/service"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func InitVideoRoutes(router *gin.RouterGroup, db *sqlx.DB) {
	videoRepo := repository.NewVideoRepository(db)
	videoService := service.NewVideoService(videoRepo)
	videoHandler := handler.NewVideoHandler(videoService)

	videoRoutes := router.Group("/videos")
	videoRoutes.Use(internal.AuthMiddleware())
	{
		videoRoutes.POST("/", videoHandler.UploadVideos())
		videoRoutes.DELETE("/:id", videoHandler.DeleteVideos())
	}
	router.GET("/videos", videoHandler.GetAllVideosWithPagination)
}
