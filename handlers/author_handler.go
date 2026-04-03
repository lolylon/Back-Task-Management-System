package handlers

import (
	"net/http"
	"strconv"

	"bookstore/config"
	"bookstore/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthorHandler struct {
	db *gorm.DB
}

func NewAuthorHandler() *AuthorHandler {
	return &AuthorHandler{
		db: config.DB,
	}
}

func (h *AuthorHandler) GetAuthors(c *gin.Context) {
	var authors []models.Author
	if err := h.db.Find(&authors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch authors"})
		return
	}

	c.JSON(http.StatusOK, authors)
}

func (h *AuthorHandler) GetAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	var author models.Author
	if err := h.db.First(&author, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch author"})
		}
		return
	}

	c.JSON(http.StatusOK, author)
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

	if err := h.db.Create(&author).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create author"})
		return
	}

	c.JSON(http.StatusCreated, author)
}

func (h *AuthorHandler) UpdateAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	var existingAuthor models.Author
	if err := h.db.First(&existingAuthor, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch author"})
		}
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

	if err := h.db.Model(&existingAuthor).Updates(updatedAuthor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update author"})
		return
	}

	h.db.First(&existingAuthor, id)

	c.JSON(http.StatusOK, existingAuthor)
}

func (h *AuthorHandler) DeleteAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	if err := h.db.Delete(&models.Author{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete author"})
		return
	}

	c.Status(http.StatusNoContent)
}