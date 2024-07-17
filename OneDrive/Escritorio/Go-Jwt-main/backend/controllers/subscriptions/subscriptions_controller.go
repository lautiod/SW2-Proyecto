package subscriptions

import (
	"backend/dao"
	domain "backend/domain/subscriptions"
	service "backend/services/subscriptions"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Subscribe(c *gin.Context) {

	var user dao.User
	val, exists := c.Get("user")
	if exists {
		user = val.(dao.User)
	}

	var request domain.SubRequest
	err := c.Bind(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to read body",
		})
		return
	}

	err = service.Subscribe(user.UserID, request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "subscribed succesfully",
	})
}
