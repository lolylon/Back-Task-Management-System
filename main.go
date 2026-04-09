package main

import (
	"bookstore/config"
	"bookstore/handlers"
	"bookstore/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()

	taskHandler := handlers.NewTaskHandler()
	userHandler := handlers.NewUserHandler()
	projectHandler := handlers.NewProjectHandler()
	authHandler := handlers.NewAuthHandler()

	r := gin.Default()

	// Public auth endpoints
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Task endpoints
		protected.GET("/tasks", taskHandler.GetTasks)
		protected.POST("/tasks", taskHandler.CreateTask)
		protected.GET("/tasks/:id", taskHandler.GetTask)
		protected.PUT("/tasks/:id", taskHandler.UpdateTask)
		protected.DELETE("/tasks/:id", taskHandler.DeleteTask)
		protected.GET("/tasks/search", taskHandler.SearchTasks)
		protected.PUT("/tasks/:id/status", taskHandler.UpdateTaskStatus)

		// User endpoints
		protected.GET("/users", userHandler.GetUsers)
		protected.POST("/users", userHandler.CreateUser)
		protected.GET("/users/:id", userHandler.GetUser)
		protected.PUT("/users/:id", userHandler.UpdateUser)
		protected.DELETE("/users/:id", userHandler.DeleteUser)

		// Project endpoints
		protected.GET("/projects", projectHandler.GetProjects)
		protected.POST("/projects", projectHandler.CreateProject)
		protected.GET("/projects/:id", projectHandler.GetProject)
		protected.PUT("/projects/:id", projectHandler.UpdateProject)
		protected.DELETE("/projects/:id", projectHandler.DeleteProject)
	}

	r.Run(":8084")
}
