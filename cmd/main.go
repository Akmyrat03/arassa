package main

import (
	_ "arassachylyk/docs"
	adminRoutes "arassachylyk/internal/admin/routes"
	catRoutes "arassachylyk/internal/categories/routes"
	contactRoutes "arassachylyk/internal/contact/routes"
	imgRoutes "arassachylyk/internal/images/routes"
	mottoRoutes "arassachylyk/internal/motto/routes"
	newsRoutes "arassachylyk/internal/news/routes"
	videoRoutes "arassachylyk/internal/videos/routes"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"arassachylyk/pkg/database"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Arassachylyk project
// @version 1.0
// @description Arassachylyk project
// @host localhost:8000
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := InitConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	DB, err := database.ConnectToDB(database.Config{
		Host:     viper.GetString("DB.host"),
		Port:     viper.GetString("DB.port"),
		Username: viper.GetString("DB.username"),
		Password: viper.GetString("DB.password"),
		DBName:   viper.GetString("DB.dbname"),
		SSLMode:  viper.GetString("DB.sslmode"),
	})

	if err != nil {
		log.Fatalf("failed to initialize DB: %v", err.Error())
	}

	app := gin.Default()

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Health check route
	app.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	corsConfig := cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://yourfrontend.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           24 * 3600,
	})

	app.Use(corsConfig)

	api := app.Group("/api")
	mottoRoutes.InitMottoRoutes(api, DB)
	catRoutes.InitCatRoutes(api, DB)
	newsRoutes.InitNewsRoutes(api, DB)
	contactRoutes.InitContactRoutes(api, DB)
	imgRoutes.InitImageRoutes(api, DB)
	videoRoutes.InitVideoRoutes(api, DB)
	adminRoutes.InitAdminRoutes(api, DB)

	// if err := app.Run(viper.GetString("APP.host")); err != nil {
	// 	log.Fatalf("Failed running app: %v", err)
	// }

	server := &http.Server{
		Addr:    viper.GetString("APP.host"),
		Handler: app,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	log.Println("Server is running...")

	//Graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down the server")

	//Timeout context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Shutdown the http server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	//Close database connection
	if err := DB.Close(); err != nil {
		log.Fatalf("Failed to close database: %v", err)
	}

	log.Println("Server exited gracefully")
}

func InitConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.SetDefault("APP.host", "localhost:8000")
	return nil
}
