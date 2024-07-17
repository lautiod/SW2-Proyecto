package comments

import (
	domain "backend/domain/comments"
	service "backend/services/comments"
	"strconv"

	_ "log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddComment(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User ID not found",
		})
		return
	}

	userID := val.(uint)

	var request domain.CommentRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	err := service.AddComment(userID, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment added successfully"})
}

func GetCommentsByCourseID(c *gin.Context) {

	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid course ID",
		})
		return
	}

	courseIDUint := uint(courseID)

	comments, err := service.GetCommentsByCourseID(courseIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, comments)
}
