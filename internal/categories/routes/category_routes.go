package routes

import (
	"arassachylyk/internal"
	"arassachylyk/internal/categories/handler"
	"arassachylyk/internal/categories/repository"
	"arassachylyk/internal/categories/service"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func InitCatRoutes(router *gin.RouterGroup, DB *sqlx.DB) {
	catRepo := repository.NewCategoryRepository(DB)
	catServ := service.NewCategoryService(catRepo)
	catHand := handler.NewCategoryHandler(catServ)

	catRoutes := router.Group("/categories")
	catRoutes.Use(internal.AuthMiddleware())
	{
		catRoutes.POST("/add", catHand.CreateCategory())
		catRoutes.DELETE("/delete/:id", catHand.DeleteCategory())
	}
	router.GET("/categories/all", catHand.GetAllCategories)

}
