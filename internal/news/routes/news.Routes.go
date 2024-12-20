package routes

import (
	"arassachylyk/internal/news/handler"
	"arassachylyk/internal/news/repository"
	"arassachylyk/internal/news/service"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func InitNewsRoutes(router *gin.RouterGroup, DB *sqlx.DB) {
	newsRepo := repository.NewRepository(DB)
	newsServ := service.NewService(newsRepo)
	newsHand := handler.NewHandler(newsServ)

	newsRoutes := router.Group("/news")
	newsRoutes.POST("/add-news", newsHand.CreateNews())
	newsRoutes.DELETE("/delete/:id", newsHand.DeleteNews())
	newsRoutes.GET("/all", newsHand.GetAllNews)
	newsRoutes.GET("/category", newsHand.GetAllNewsByLangAndCategory)
}
