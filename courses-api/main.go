package main

import (
	"courses-api/clients/queues"
	controllers "courses-api/controllers/courses"
	repositories "courses-api/repositories/courses"
	services "courses-api/services/courses"

	"log"

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
	service := services.NewService(mainRepository, eventsQueue)

	// Controllers
	controller := controllers.NewController(service)

	// Router
	router := gin.Default()
	router.GET("/courses", controller.GetCourses)
	router.GET("/courses/:id", controller.GetCourseByID)
	router.POST("/courses", controller.CreateCourse)
	router.PUT("/courses/:id", controller.UpdateCourse)
	//inscr

	if err := router.Run(":8081"); err != nil {
		log.Fatalf("error running application: %v", err)
	}

}
