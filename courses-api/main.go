package main

import (
	"courses-api/clients/queues"
	controllers "courses-api/controllers/courses"
	repositories "courses-api/repositories/courses"
	services "courses-api/services/courses"

	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// MongoDB
	mainRepository := repositories.NewMongo(repositories.MongoConfig{
		Host:       "mongo", // Colocar 'mongo' para correr con docker
		Port:       "27017",
		Username:   "root",
		Password:   "root",
		Database:   "courses-api",
		Collection: "courses",
	})

	// Rabbit
	eventsQueue := queues.NewRabbit(queues.RabbitConfig{
		Host:      "rabbitmq",
		Port:      "5672",
		Username:  "root",
		Password:  "root",
		QueueName: "courses-news",
	})

	// Services
	courseService := services.NewService(mainRepository, eventsQueue)
	// Controllers
	courseController := controllers.NewController(courseService)

	// Router
	router := gin.Default()

	// Cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/courses", courseController.GetCourses)
	router.GET("/courses/:id", courseController.GetCourseByID)
	router.POST("/courses", courseController.CreateCourse)
	router.PUT("/courses/:id", courseController.UpdateCourse)

	if err := router.Run(":8081"); err != nil {
		log.Fatalf("error running application: %v", err)
	}

}
