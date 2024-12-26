package routes

import (
	"arassachylyk/internal"
	"arassachylyk/internal/news/handler"
	"arassachylyk/internal/news/repository"
	"arassachylyk/internal/news/service"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func InitNewsRoutes(router *gin.RouterGroup, db *sqlx.DB) {
	newsRepo := repository.NewRepository(db)
	newsServ := service.NewService(newsRepo)
	newsHand := handler.NewHandler(newsServ)

	newsRoutes := router.Group("/news")
	newsRoutes.Use(internal.AuthMiddleware())
	{
		newsRoutes.DELETE("/:id", newsHand.DeleteNews())
		newsRoutes.POST("/", newsHand.CreateNews())
	}
	router.GET("/news/all", newsHand.GetAllNewsPagination)
	router.GET("/news", newsHand.GetAllNewsByLangAndCategory)
}
