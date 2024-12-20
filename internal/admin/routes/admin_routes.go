package routes

import (
	"arassachylyk/internal/admin/middleware"
	"arassachylyk/internal/admin/repository"
	"arassachylyk/internal/admin/service"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func InitAdminRoutes(router *gin.RouterGroup, DB *sqlx.DB) {
	adminRepo := repository.NewAdminRepository(DB)
	adminServ := service.NewAdminService(adminRepo)
	adminMidd := middleware.NewAdminMiddleware(adminRepo, adminServ)

	adminRoutes := router.Group("/admin")
	{
		adminRoutes.POST("/signup", adminMidd.SignUp())
		adminRoutes.POST("/login", adminMidd.Login())
		adminRoutes.GET("/profile", adminMidd.Profile())
	}

}
