package main

import (
	"backend/clients"
	"backend/controllers/comments"
	"backend/controllers/courses"
	"backend/controllers/subscriptions"
	"backend/controllers/users"
	"backend/middleware"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	clients.LoadEnvVariables()
	clients.ConnectToDb()
	clients.SyncDb()
}

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// users
	r.POST("/signup", users.Signup)
	r.POST("/login", users.Login)
	r.GET("/myinfo", middleware.RequireAuth, users.MyInfo)
	r.GET("/mycourses", middleware.RequireAuth, users.GetCoursesByUserID)
	// courses
	r.GET("/courses", courses.GetCourses)
	r.GET("/courses/:id", courses.GetCourseById)
	r.POST("/createcourse", middleware.RequireAuth, courses.CreateCourse)
	r.POST("/mycourse", middleware.RequireAuth, courses.CourseSubs)
	r.GET("/courses/search", middleware.RequireAuth, courses.GetCoursesBySearch)
	r.POST("/addfile", middleware.RequireAuth, courses.AddFile)
	// subscriptions
	r.POST("/subscribe", middleware.RequireAuth, subscriptions.Subscribe)
	// comments
	r.POST("/comment", middleware.RequireAuth, comments.AddComment)
	r.GET("/comment/:id", middleware.RequireAuth, comments.GetCommentsByCourseID)

	r.Run()
}
