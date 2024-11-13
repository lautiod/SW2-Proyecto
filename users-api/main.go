package main

import (
	"log"
	"time"
	controllers "users-api/controllers/users"
	middleware "users-api/middleware"
	repositories "users-api/repositories/users"
	services "users-api/services/users"

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

	// Services
	service := services.NewService(mySQLRepo, cacheRepo, memcachedRepo)

	// Handlers
	controller := controllers.NewController(service)

	// Create router
	router := gin.Default()

	// URL mappings
	router.GET("/users", middleware.RequireAuth, controller.GetAll)
	router.GET("/users/:id", middleware.RequireAuth, controller.GetByID)
	router.POST("/users", controller.Create)
	router.PUT("/users/:id", middleware.RequireAuth, controller.Update)
	router.POST("/login", controller.Login)

	// Run application
	if err := router.Run(":8080"); err != nil {
		log.Panicf("Error running application: %v", err)
	}
}
