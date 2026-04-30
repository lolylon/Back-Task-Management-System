package handlers

import (
	"net/http"
	"strconv"

	"bookstore/config"
	"bookstore/models"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

type CommentHandler struct {
	db     *gorm.DB
	client *resty.Client
}

func NewCommentHandler() *CommentHandler {
	client := resty.New()
	client.SetTimeout(10 * 1000) // 10 seconds timeout

	return &CommentHandler{
		db:     config.DB,
		client: client,
	}
}

func (h *CommentHandler) GetComments(c *gin.Context) {
	taskIDStr := c.Query("task_id")
	var comments []models.Comment

	query := h.db
	if taskIDStr != "" {
		taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}
		query = query.Where("task_id = ?", taskID)
	}

	if err := query.Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *CommentHandler) GetComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var comment models.Comment
	if err := h.db.First(&comment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comment"})
		}
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	var comment models.Comment

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if comment.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content is required"})
		return
	}

	if comment.TaskID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	if comment.UserEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User email is required"})
		return
	}

	if err := h.db.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Use restyV2 to call an external notification service
	h.sendNotification(comment)

	c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandler) UpdateComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var existingComment models.Comment
	if err := h.db.First(&existingComment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comment"})
		}
		return
	}

	var updatedComment models.Comment
	if err := c.ShouldBindJSON(&updatedComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if updatedComment.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content is required"})
		return
	}

	if err := h.db.Model(&existingComment).Updates(updatedComment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	h.db.First(&existingComment, id)
	c.JSON(http.StatusOK, existingComment)
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	if err := h.db.Delete(&models.Comment{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CommentHandler) sendNotification(comment models.Comment) {
	// Example of using restyV2 to call an external API
	// This calls a test API (JSONPlaceholder) to demonstrate restyV2 usage
	resp, err := h.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"comment_id": comment.ID,
			"task_id":    comment.TaskID,
			"user_email": comment.UserEmail,
			"content":    comment.Content,
			"created_at": comment.CreatedAt,
		}).
		Post("https://jsonplaceholder.typicode.com/posts")

	if err != nil {
		// Log error but don't fail the request
		return
	}

	// Log response for debugging - shows restyV2 is working
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		// Successfully called external API using restyV2
	}
}

func (h *CommentHandler) TestExternalAPI(c *gin.Context) {
	// Explicit endpoint to demonstrate restyV2 usage
	resp, err := h.client.R().
		SetHeader("Accept", "application/json").
		Get("https://jsonplaceholder.typicode.com/posts/1")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call external API", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode(),
		"body":        resp.String(),
		"message":     "Called using restyV2",
	})
}
