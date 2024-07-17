package users

import (
	"backend/clients"
	_ "backend/dao"
	domain "backend/domain/users"
	service "backend/services/users"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {

	var request domain.Signup_Request

	err := c.Bind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}
	log.Printf("Received User: %+v", request)

	err = service.Signup(request.Email, request.Password, request.IsAdmin)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User added successfully"})
}

func Login(c *gin.Context) {
	var body domain.Login_Request

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	response, tokenString, err := service.Login(body)

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

func MyInfo(c *gin.Context) {

	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User ID not found",
		})
		return
	}

	userID := val.(uint)

	courseIDs, err := clients.GetCourseIDsByUserID(userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    userID,
		"courses": courseIDs,
	})
}

func GetCoursesByUserID(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User ID not found",
		})
		return
	}

	userID := val.(uint)

	courses, err := service.GetCoursesByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, courses)
}
