package courses

import (
	"backend/clients"
	"backend/dao"
	domainC "backend/domain/courses"
	domainSub "backend/domain/subscriptions"
	service "backend/services/courses"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddFile(c *gin.Context) {

	// val, exists := c.Get("userID")
	// if !exists {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"message": "User not found",
	// 	})
	// 	return
	// }

	// userID := val.(uint)

	var request domainC.FileRequest
	err := c.Bind(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to read body",
		})
		return
	}

	// err = service.CheckUsers(userID, request.CourseID)
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"message": err.Error(),
	// 	})
	// 	return
	// }

	err = service.AddFile(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
}

func CreateCourse(c *gin.Context) {

	val, exists := c.Get("isAdmin")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User ID not found",
		})
		return
	}

	admin := val.(bool)

	if !admin {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "user is not admin",
		})
		return
	}

	var course domainC.CourseDetail

	err := c.Bind(&course)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to read body",
		})
		return
	}

	val, exists = c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User ID not found",
		})
		return
	}

	userID := val.(uint)

	err = service.CreateCourse(course, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to read body",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Course created successfully"})

}

func GetCourseById(c *gin.Context) {
	var course domainC.CourseDetail
	id, _ := strconv.Atoi(c.Param("id"))

	course, err := service.GetCourseById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, course)
}

func GetCourses(c *gin.Context) {
	var courses domainC.CoursesDetail

	courses, err := service.GetCourses()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, courses)
}

func GetCoursesBySearch(c *gin.Context) {
	query := c.Query("q")

	courses, err := service.GetCoursesBySearch(query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to search courses: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, courses)
}

func CourseSubs(c *gin.Context) {
	//Saco el user mediante la galletita :)
	var user dao.User
	val, exists := c.Get("user")
	if exists {
		user = val.(dao.User)
	}

	//bindeo request para despues recuperar el curso completo :p
	var request domainSub.SubRequest
	err := c.Bind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to read body",
		})
		return
	}

	//recupero el curso pasandole el id de curso que estaba en la request :3
	var course dao.Course
	course, err = clients.SelectCourseByID(request.CourseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	//valido que el creatorID sea el mismo que el userID para que no haya problemas uWu
	if user.UserID != course.CreatorID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": errors.New("que haces desubicado"),
		})
		return
	}

	//una vez validada la condicion saco el array de usersIDs >:/
	userIDs, err := clients.GetUserIDsByCourseID(course.CourseID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": course,
		"users":   userIDs,
	})
}
