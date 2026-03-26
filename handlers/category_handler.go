package handlers

import (
	"net/http"
	"strconv"

	"bookstore/models"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categories []models.Category
	nextID     int
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		categories: []models.Category{},
		nextID:     1,
	}
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	c.JSON(http.StatusOK, h.categories)
}

func (h *CategoryHandler) GetCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	for _, category := range h.categories {
		if category.ID == id {
			c.JSON(http.StatusOK, category)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var category models.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if category.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	category.ID = h.nextID
	h.nextID++
	h.categories = append(h.categories, category)

	c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var updatedCategory models.Category
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if updatedCategory.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	for i, category := range h.categories {
		if category.ID == id {
			updatedCategory.ID = id
			h.categories[i] = updatedCategory
			c.JSON(http.StatusOK, updatedCategory)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	for i, category := range h.categories {
		if category.ID == id {
			h.categories = append(h.categories[:i], h.categories[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
}