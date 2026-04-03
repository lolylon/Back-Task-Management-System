package handlers

import (
	"net/http"
	"strconv"

	"bookstore/config"
	"bookstore/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookHandler struct {
	db *gorm.DB
}

func NewBookHandler() *BookHandler {
	return &BookHandler{
		db: config.DB,
	}
}

func (h *BookHandler) GetBooks(c *gin.Context) {
	page := 1
	limit := 10
	sort := c.Query("sort")
	search := c.Query("search")

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset := (page - 1) * limit

	var books []models.Book
	var total int64

	query := h.db.Model(&models.Book{}).Preload("Author").Preload("Category")

	if search != "" {
		query = query.Where("title ILIKE ? OR isbn ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if sort == "price_asc" {
		query = query.Order("price asc")
	} else if sort == "price_desc" {
		query = query.Order("price desc")
	} else if sort == "title" {
		query = query.Order("title asc")
	} else {
		query = query.Order("created_at desc")
	}

	if err := query.Offset(offset).Limit(limit).Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}

	query.Count(&total)

	response := gin.H{
		"books": books,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}

	c.JSON(http.StatusOK, response)
}

func (h *BookHandler) GetBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var book models.Book
	if err := h.db.Preload("Author").Preload("Category").First(&book, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch book"})
		}
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var book models.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if book.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	if book.Price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price cannot be negative"})
		return
	}

	if book.AuthorID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Author ID is required"})
		return
	}

	if book.CategoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category ID is required"})
		return
	}

	if err := h.db.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	h.db.Preload("Author").Preload("Category").First(&book, book.ID)

	c.JSON(http.StatusCreated, book)
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var existingBook models.Book
	if err := h.db.First(&existingBook, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch book"})
		}
		return
	}

	var updatedBook models.Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if updatedBook.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	if updatedBook.Price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price cannot be negative"})
		return
	}

	if updatedBook.AuthorID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Author ID is required"})
		return
	}

	if updatedBook.CategoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category ID is required"})
		return
	}

	if err := h.db.Model(&existingBook).Updates(updatedBook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	h.db.Preload("Author").Preload("Category").First(&existingBook, id)

	c.JSON(http.StatusOK, existingBook)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	if err := h.db.Delete(&models.Book{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *BookHandler) SearchBooks(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	var books []models.Book
	if err := h.db.Preload("Author").Preload("Category").
		Where("title ILIKE ? OR isbn ILIKE ? OR EXISTS (SELECT 1 FROM authors WHERE authors.id = books.author_id AND authors.name ILIKE ?) OR EXISTS (SELECT 1 FROM categories WHERE categories.id = books.category_id AND categories.name ILIKE ?)", 
			"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%").
		Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search books"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"books": books,
		"query": query,
		"count": len(books),
	})
}