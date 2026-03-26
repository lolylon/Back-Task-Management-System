package handlers

import (
	"net/http"
	"strconv"

	"bookstore/models"

	"github.com/gin-gonic/gin"
)

type AuthorHandler struct {
	authors []models.Author
	nextID  int
}

func NewAuthorHandler() *AuthorHandler {
	return &AuthorHandler{
		authors: []models.Author{},
		nextID:  1,
	}
}

func (h *AuthorHandler) GetAuthors(c *gin.Context) {
	c.JSON(http.StatusOK, h.authors)
}

func (h *AuthorHandler) GetAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	for _, author := range h.authors {
		if author.ID == id {
			c.JSON(http.StatusOK, author)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
}

func (h *AuthorHandler) CreateAuthor(c *gin.Context) {
	var author models.Author

	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if author.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	author.ID = h.nextID
	h.nextID++
	h.authors = append(h.authors, author)

	c.JSON(http.StatusCreated, author)
}

func (h *AuthorHandler) UpdateAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	var updatedAuthor models.Author
	if err := c.ShouldBindJSON(&updatedAuthor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if updatedAuthor.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	for i, author := range h.authors {
		if author.ID == id {
			updatedAuthor.ID = id
			h.authors[i] = updatedAuthor
			c.JSON(http.StatusOK, updatedAuthor)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
}

func (h *AuthorHandler) DeleteAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	for i, author := range h.authors {
		if author.ID == id {
			h.authors = append(h.authors[:i], h.authors[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
}