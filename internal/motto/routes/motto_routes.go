package routes

import (
	"arassachylyk/internal"
	"arassachylyk/internal/motto/handler"
	"arassachylyk/internal/motto/repository"
	"arassachylyk/internal/motto/service"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func InitMottoRoutes(router *gin.RouterGroup, db *sqlx.DB) {
	mottoRepo := repository.NewYearRepository(db)
	mottoService := service.NewYearService(mottoRepo)
	mottoHandler := handler.NewYearHandler(mottoService)

	mottoRoutes := router.Group("/motto")
	mottoRoutes.Use(internal.AuthMiddleware())
	{
		mottoRoutes.POST("/", mottoHandler.AddMotto())
		mottoRoutes.DELETE("/:id", mottoHandler.DeleteMotto())
	}
	mottoRoutes.GET("/", mottoHandler.GetAllMottos())

}
