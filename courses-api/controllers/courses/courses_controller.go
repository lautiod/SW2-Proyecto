package courses

import (
	"context"
	coursesDomain "courses-api/domain/courses"
	inscriptionsDomain "courses-api/domain/inscriptions"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Service interface {
	GetCourses(ctx context.Context) (coursesDomain.Courses, error)
	GetCourseByID(ctx context.Context, id string) (coursesDomain.Course, error)
	CreateCourse(ctx context.Context, course coursesDomain.Course) (string, error)
	UpdateCourse(ctx context.Context, course coursesDomain.Course) error
	EnrollUser(ctx context.Context, inscription inscriptionsDomain.Inscription) (string, error)
	GetCoursesByUserID(ctx context.Context, userID string) ([]coursesDomain.Course, error)
}

type Controller struct {
	service Service
}

func NewController(service Service) Controller {
	return Controller{
		service: service,
	}
}

func (controller Controller) GetCourses(ctx *gin.Context) {
	courses, err := controller.service.GetCourses(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error getting courses: %s", err.Error()),
		})
		return
	}
	ctx.JSON(http.StatusOK, courses)
}

func (controller Controller) GetCourseByID(ctx *gin.Context) {
	// Validate ID param
	courseID := strings.TrimSpace(ctx.Param("id"))

	// Get hotel by ID using the service
	course, err := controller.service.GetCourseByID(ctx.Request.Context(), courseID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error getting course: %s", err.Error()),
		})
		return
	}

	// Send response
	ctx.JSON(http.StatusOK, course)
}

func (controller Controller) CreateCourse(ctx *gin.Context) {
	// Parse course
	var course coursesDomain.Course
	if err := ctx.ShouldBindJSON(&course); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	// Create hotel
	id, err := controller.service.CreateCourse(ctx.Request.Context(), course)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error creating course: %s", err.Error()),
		})
		return
	}

	// Send ID
	ctx.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func (controller Controller) UpdateCourse(ctx *gin.Context) {
	// Validate ID param
	id := strings.TrimSpace(ctx.Param("id"))

	// Parse hotel
	var course coursesDomain.Course
	if err := ctx.ShouldBindJSON(&course); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	// Set the ID from the URL to the hotel object
	course.ID = id

	// Update hotel
	if err := controller.service.UpdateCourse(ctx.Request.Context(), course); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error updating course: %s", err.Error()),
		})
		return
	}

	// Send response
	ctx.JSON(http.StatusOK, gin.H{
		"message": id,
	})
}

func (controller Controller) EnrollUser(ctx *gin.Context) {
	// Parse inscription
	var inscription inscriptionsDomain.Inscription
	if err := ctx.ShouldBindJSON(&inscription); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	// Create hotel
	id, err := controller.service.EnrollUser(ctx.Request.Context(), inscription)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error enrrolling user: %s", err.Error()),
		})
		return
	}

	// Send ID
	ctx.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func (controller Controller) GetCoursesByUserID(ctx *gin.Context) {

	userID := strings.TrimSpace(ctx.Param("id"))

	courses, err := controller.service.GetCoursesByUserID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error getting courses: %s", err.Error()),
		})
		return
	}
	ctx.JSON(http.StatusOK, courses)
}
