package main

import (
	"bookstore/config"
	"bookstore/handlers"

	"github.com/gin-gonic/gin"
)


func main() {
	config.InitDB()

	bookHandler := handlers.NewBookHandler()
	authorHandler := handlers.NewAuthorHandler()
	categoryHandler := handlers.NewCategoryHandler()

	r := gin.Default()

	r.GET("/books", bookHandler.GetBooks)
	r.POST("/books", bookHandler.CreateBook)
	r.GET("/books/:id", bookHandler.GetBook)
	r.PUT("/books/:id", bookHandler.UpdateBook)
	r.DELETE("/books/:id", bookHandler.DeleteBook)
	r.GET("/books/search", bookHandler.SearchBooks)

	r.GET("/authors", authorHandler.GetAuthors)
	r.POST("/authors", authorHandler.CreateAuthor)
	r.GET("/authors/:id", authorHandler.GetAuthor)
	r.PUT("/authors/:id", authorHandler.UpdateAuthor)
	r.DELETE("/authors/:id", authorHandler.DeleteAuthor)

	r.GET("/categories", categoryHandler.GetCategories)
	r.POST("/categories", categoryHandler.CreateCategory)
	r.GET("/categories/:id", categoryHandler.GetCategory)
	r.PUT("/categories/:id", categoryHandler.UpdateCategory)
	r.DELETE("/categories/:id", categoryHandler.DeleteCategory)

	r.Run(":8080")
}