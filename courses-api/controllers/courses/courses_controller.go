package courses

import (
	"context"
	coursesDomain "courses-api/domain/courses"
	inscriptionsDomain "courses-api/domain/inscriptions"
	"fmt"
	"log"
	"net/http"
	"strings"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Service interface {
	GetCourses(ctx context.Context) (coursesDomain.Courses, error)
	GetCourseByID(ctx context.Context, id string) (coursesDomain.Course, error)
	CreateCourse(ctx context.Context, course coursesDomain.Course, userID string) (string, error)
	UpdateCourse(ctx context.Context, course coursesDomain.Course, userID string) error
	EnrollUser(ctx context.Context, inscription inscriptionsDomain.Inscription) (string, error)
	GetCoursesByUserID(ctx context.Context, userID string) ([]coursesDomain.Course, error)
	GetCoursesDisponibility(ctx context.Context) (coursesDomain.Courses, error)
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
	// Obtener el token de la cookie
	tokenString, err := ctx.Cookie("Authorization")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "no authorization cookie found",
		})
		return
	}

	// Extraer el userID del token
	userID, err := extractUserIDFromToken(tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	log.Printf("Token extracted userID: %s", userID) // Log para debug

	var course coursesDomain.Course
	if err := ctx.ShouldBindJSON(&course); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	id, err := controller.service.CreateCourse(ctx.Request.Context(), course, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "unauthorized") {
			statusCode = http.StatusUnauthorized
		}
		ctx.JSON(statusCode, gin.H{
			"error": fmt.Sprintf("error creating course: %s", err.Error()),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func extractUserIDFromToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return "", fmt.Errorf("error parsing token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok := claims["sub"].(float64); ok {
			return fmt.Sprintf("%d", int64(sub)), nil
		}
	}

	return "", fmt.Errorf("invalid token claims")
}

func (controller Controller) UpdateCourse(ctx *gin.Context) {
	// Obtener el token de la cookie
	tokenString, err := ctx.Cookie("Authorization")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "no authorization cookie found",
		})
		return
	}

	// Extraer el userID del token
	userID, err := extractUserIDFromToken(tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	log.Printf("Token extracted userID: %s", userID) // Log para debug

	id := strings.TrimSpace(ctx.Param("id"))
	var course coursesDomain.Course
	if err := ctx.ShouldBindJSON(&course); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	course.ID = id
	if err := controller.service.UpdateCourse(ctx.Request.Context(), course, userID); err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "unauthorized") {
			statusCode = http.StatusUnauthorized
		}
		ctx.JSON(statusCode, gin.H{
			"error": fmt.Sprintf("error updating course: %s", err.Error()),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": id,
	})
}

func (controller Controller) EnrollUser(ctx *gin.Context) {

	var inscriptionRequest inscriptionsDomain.Inscription
	if err := ctx.ShouldBindJSON(&inscriptionRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	courseID := strings.TrimSpace(ctx.Param("id"))

	inscription := inscriptionsDomain.Inscription{
		CourseID: courseID,
		UserID:   inscriptionRequest.UserID,
	}

	// Enroll User
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

func (controller Controller) GetCoursesDisponibility(ctx *gin.Context) {
	availablesCourses, err := controller.service.GetCoursesDisponibility(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error getting courses disponibility: %s", err.Error()),
		})
		return
	}
	ctx.JSON(http.StatusOK, availablesCourses)
}
