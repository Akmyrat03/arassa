package main

import (
	_ "arassachylyk/docs"
	catRoutes "arassachylyk/internal/categories/routes"
	yearRoutes "arassachylyk/internal/motto/routes"
	newsRoutes "arassachylyk/internal/news/routes"
	"arassachylyk/pkg/database"
	"log"

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

	// Add CORS middleware
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:3000"}, // Frontend URL
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// }))

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := app.Group("/api")

	yearRoutes.InitYearRoutes(api, DB)
	catRoutes.InitCatRoutes(api, DB)
	newsRoutes.InitNewsRoutes(api, DB)

	if err := app.Run(viper.GetString("APP.host")); err != nil {
		log.Fatalf("Failed running app: %v", err)
	}
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
