package routes

import (
	"arassachylyk/internal"
	"arassachylyk/internal/categories/handler"
	"arassachylyk/internal/categories/repository"
	"arassachylyk/internal/categories/service"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func InitCatRoutes(router *gin.RouterGroup, db *sqlx.DB) {
	catRepo := repository.NewCategoryRepository(db)
	catServ := service.NewCategoryService(catRepo)
	catHand := handler.NewCategoryHandler(catServ)

	catRoutes := router.Group("/categories")
	catRoutes.Use(internal.AuthMiddleware())
	{
		catRoutes.POST("/", catHand.CreateCategory())
		catRoutes.DELETE("/:id", catHand.DeleteCategory())
	}
	router.GET("/categories", catHand.GetAllCategories())
}
