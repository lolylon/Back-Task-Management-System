package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"bookstore/models"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	books  []models.Book
	nextID int
}

func NewBookHandler() *BookHandler {
	return &BookHandler{
		books:  []models.Book{},
		nextID: 1,
	}
}

func (h *BookHandler) GetBooks(c *gin.Context) {
	page := 1
	limit := 10
	category := ""

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

	if cParam := c.Query("category"); cParam != "" {
		category = strings.ToLower(cParam)
	}

	filtered := h.books

	if category != "" {
		filtered = []models.Book{}
		for _, book := range h.books {
			if strings.Contains(strings.ToLower(book.Title), category) {
				filtered = append(filtered, book)
			}
		}
	}

	total := len(filtered)
	start := (page - 1) * limit
	end := start + limit

	if start >= total {
		filtered = []models.Book{}
	} else if end > total {
		filtered = filtered[start:total]
	} else {
		filtered = filtered[start:end]
	}

	response := gin.H{
		"books": filtered,
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
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	for _, book := range h.books {
		if book.ID == id {
			c.JSON(http.StatusOK, book)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
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

	book.ID = h.nextID
	h.nextID++
	h.books = append(h.books, book)

	c.JSON(http.StatusCreated, book)
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
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

	for i, book := range h.books {
		if book.ID == id {
			updatedBook.ID = id
			h.books[i] = updatedBook
			c.JSON(http.StatusOK, updatedBook)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	for i, book := range h.books {
		if book.ID == id {
			h.books = append(h.books[:i], h.books[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
}