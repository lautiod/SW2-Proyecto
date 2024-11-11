package main

import (
	controllers "inscriptions-api/controllers/inscriptions"
	repositories "inscriptions-api/repositories/inscriptions"
	services "inscriptions-api/services/inscriptions"

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
		Database:   "inscriptions-api",
		Collection: "inscriptions",
	})

	// Services
	service := services.NewService(mainRepository)
	// Controllers
	controller := controllers.NewController(service)

	// Router
	router := gin.Default()
	router.POST("/inscriptions/courses", controller.EnrollUser)

	if err := router.Run(":8083"); err != nil {
		log.Fatalf("error running application: %v", err)
	}

}
