package main

import (
	"log"
	"time"
	controllers "users-api/controllers/users"
	"users-api/middleware"
	repositories "users-api/repositories/users"
	services "users-api/services/users"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// MySQL
	mySQLRepo := repositories.NewMySQL(
		repositories.MySQLConfig{
			Host:     "mysql",
			Port:     "3306",
			Database: "users-api",
			Username: "root",
			Password: "root",
		},
	)

	// Cache
	cacheRepo := repositories.NewCache(repositories.CacheConfig{
		TTL: 1 * time.Minute,
	})

	// // Memcached
	memcachedRepo := repositories.NewMemcached(repositories.MemcachedConfig{
		Host: "memcached",
		Port: "11211",
	})

	// // Tokenizer
	// jwtTokenizer := tokenizers.NewTokenizer(
	// 	tokenizers.JWTConfig{
	// 		Key:      "ThisIsAnExampleJWTKey!",
	// 		Duration: 1 * time.Hour,
	// 	},
	// )

	// // Services
	service := services.NewService(mySQLRepo, cacheRepo, memcachedRepo /*, jwtTokenizer*/)

	// Handlers
	controller := controllers.NewController(service)

	// Create router
	router := gin.Default()
	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// URL mappings
	router.GET("/users", controller.GetAll)
	router.GET("/users/:id", controller.GetByID)
	router.POST("/users", controller.Create)
	router.PUT("/users/:id", controller.Update)
	router.POST("/login", controller.Login)
	router.GET("/users/admin", middleware.RequireAuth, controller.GetServices)
	// Run application
	if err := router.Run(":8080"); err != nil {
		log.Panicf("Error running application: %v", err)
	}
}
