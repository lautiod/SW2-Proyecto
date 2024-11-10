package main

import (
	"log"
	"search-api/clients/queues"
	controllers "search-api/controllers/search"
	repositories "search-api/repositories/courses"
	services "search-api/services/search"

	"github.com/gin-gonic/gin"
)

func main() {
	// Solr
	solrRepo := repositories.NewSolr(repositories.SolrConfig{
		Host:       "solr",    // Solr host
		Port:       "8983",    // Solr port
		Collection: "courses", // Collection name
	})

	// Rabbit
	eventsQueue := queues.NewRabbit(queues.RabbitConfig{
		Host:      "rabbitmq",
		Port:      "5672",
		Username:  "root",
		Password:  "root",
		QueueName: "courses-news",
	})

	// Courses API
	coursesAPI := repositories.NewHTTP(repositories.HTTPConfig{
		Host: "courses-api",
		Port: "8081",
	})

	// Services
	service := services.NewService(solrRepo, coursesAPI)

	// Controllers
	controller := controllers.NewController(service)

	// Launch rabbit consumer
	if err := eventsQueue.StartConsumer(service.HandleCourseNew); err != nil {
		log.Fatalf("Error running consumer: %v", err)
	}

	// Create router
	router := gin.Default()
	router.GET("/search", controller.Search)
	if err := router.Run(":8082"); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
