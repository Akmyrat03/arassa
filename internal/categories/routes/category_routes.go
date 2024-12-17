package routes

import (
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
	catRoutes.POST("/add", catHand.CreateCategory)
	catRoutes.DELETE("/delete/:id", catHand.DeleteCategory)
	// catRoutes.PUT("/update/:id", catHand.UpdateCategory)
	catRoutes.GET("/tkm", catHand.GetAllCategoriesTKM)
	catRoutes.GET("/eng", catHand.GetAllCategoriesENG)
	catRoutes.GET("/rus", catHand.GetAllCategoriesRUS)

}
