package users

import (
	"fmt"
	"net/http"
	"strconv"
	domain "users-api/domain/users"

	"github.com/gin-gonic/gin"
)

type Service interface {
	GetAll() ([]domain.User, error)
	GetByID(id int64) (domain.User, error)
	Create(user domain.User) (int64, error)
	Update(user domain.User) error
	Delete(id int64) error
	Login(body domain.Login_Request) (domain.LoginResponse, string, error)
}

type Controller struct {
	service Service
}

func NewController(service Service) Controller {
	return Controller{
		service: service,
	}
}

func (controller Controller) GetAll(c *gin.Context) {
	// Invoke service
	users, err := controller.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error getting all users: %s", err.Error()),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, users)
}

func (controller Controller) GetByID(c *gin.Context) {
	// Parse user ID from HTTP request
	userID := c.Param("id")
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	// Invoke service
	user, err := controller.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("user not found: %s", err.Error()),
		})
		return
	}

	// Send user
	c.JSON(http.StatusOK, user)
}

func (controller Controller) Create(c *gin.Context) {
	// Parse user from HTTP Request
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	// Invoke service
	id, err := controller.service.Create(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error creating user: %s", err.Error()),
		})
		return
	}

	// Send ID
	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func (controller Controller) Update(c *gin.Context) {
	// Parse user ID from HTTP request
	userID := c.Param("id")
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	// Parse updated user data from HTTP request
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	// Set the ID of the user to be updated
	user.ID = id

	// Invoke service
	if err := controller.service.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error updating user: %s", err.Error()),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, user)
}

func (controller Controller) Delete(c *gin.Context) {
	// Parse user ID from HTTP request
	userID := c.Param("id")
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %s", err.Error()),
		})
		return
	}

	// Invoke service
	if err := controller.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error deleting user: %s", err.Error()),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (controller Controller) Login(c *gin.Context) {
	var body domain.Login_Request

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	response, tokenString, err := controller.service.Login(body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, response)
}
