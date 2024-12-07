package routes

import (
	"arassachylyk/internal/motto/handler"
	"arassachylyk/internal/motto/repository"
	"arassachylyk/internal/motto/service"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func InitYearRoutes(router *gin.RouterGroup, db *sqlx.DB) {
	mottoRepo := repository.NewYearRepository(db)
	mottoService := service.NewYearService(mottoRepo)
	mottoHandler := handler.NewYearHandler(mottoService)

	mottoRoutes := router.Group("/motto")

	mottoRoutes.POST("/add", mottoHandler.AddYear)
	mottoRoutes.DELETE("/delete/:id", mottoHandler.DeleteYear)

}
