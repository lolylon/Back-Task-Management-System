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
	commentHandler := handlers.NewCommentHandler()

	r := gin.Default()

	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/tasks", taskHandler.GetTasks)
		protected.POST("/tasks", taskHandler.CreateTask)
		protected.GET("/tasks/:id", taskHandler.GetTask)
		protected.PUT("/tasks/:id", taskHandler.UpdateTask)
		protected.DELETE("/tasks/:id", taskHandler.DeleteTask)
		protected.GET("/tasks/search", taskHandler.SearchTasks)
		protected.PUT("/tasks/:id/status", taskHandler.UpdateTaskStatus)

		protected.GET("/users", userHandler.GetUsers)
		protected.POST("/users", userHandler.CreateUser)
		protected.GET("/users/:id", userHandler.GetUser)
		protected.PUT("/users/:id", userHandler.UpdateUser)
		protected.DELETE("/users/:id", userHandler.DeleteUser)

		protected.GET("/projects", projectHandler.GetProjects)
		protected.POST("/projects", projectHandler.CreateProject)
		protected.GET("/projects/:id", projectHandler.GetProject)
		protected.PUT("/projects/:id", projectHandler.UpdateProject)
		protected.DELETE("/projects/:id", projectHandler.DeleteProject)

		protected.GET("/comments", commentHandler.GetComments)
		protected.POST("/comments", commentHandler.CreateComment)
		protected.GET("/comments/:id", commentHandler.GetComment)
		protected.PUT("/comments/:id", commentHandler.UpdateComment)
		protected.DELETE("/comments/:id", commentHandler.DeleteComment)
		protected.GET("/comments/test-external-api", commentHandler.TestExternalAPI)
	}

	r.Run(":8085")
}
