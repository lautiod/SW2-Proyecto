package users

import (
	"context"
	"fmt"
	"net/http"
	coursesDomain "search-api/domain/courses"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Search(ctx context.Context, query string, offset int, limit int) ([]coursesDomain.Course, error)
}

type Controller struct {
	service Service
}

func NewController(service Service) Controller {
	return Controller{
		service: service,
	}
}

func (controller Controller) Search(c *gin.Context) {
	// Parse query from URL
	query := c.Query("q")

	// Parse offset from URL
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errRRor": fmt.Sprintf("invalid request: %s", err),
		})
		return
	}

	// Parse limit from URL
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erRror": fmt.Sprintf("invalid request: %s", err),
		})
		return
	}

	fmt.Printf("URI recibida: %s\n", c.Request.RequestURI)
	fmt.Printf("offset: %s, limit: %s , query: %s\n", c.Query("offset"), c.Query("limit"), c.Query("q"))

	// Invoke service
	courses, err := controller.service.Search(c.Request.Context(), query, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error searching courses: %s", err.Error()),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, courses)
}
