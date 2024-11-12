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
		Host:                   "mongo", // Colocar 'mongo' para correr con docker
		Port:                   "27017",
		Username:               "root",
		Password:               "root",
		Database:               "courses-api",
		CoursesCollection:      "courses",
		InscriptionsCollection: "inscriptions",
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
	service := services.NewService(mainRepository, eventsQueue)
	// Controllers
	controller := controllers.NewController(service)

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

	router.GET("/courses", controller.GetCourses)
	router.GET("/courses/:id", controller.GetCourseByID)
	router.POST("/courses", controller.CreateCourse)
	router.PUT("/courses/:id", controller.UpdateCourse)
	router.POST("/inscriptions/courses", controller.EnrollUser)
	router.GET("/inscriptions/user/:id", controller.GetCoursesByUserID)

	if err := router.Run(":8081"); err != nil {
		log.Fatalf("error running application: %v", err)
	}

}
