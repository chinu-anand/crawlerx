package api

import (
	"net/http"

	"github.com/chinu-anand/crawlerx/internal/models"
	"github.com/chinu-anand/crawlerx/internal/queue"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var QueueService queue.JobQueue

func PostCrawl(c *gin.Context) {
	var req CrawlRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id := uuid.New().String()

	job := models.CrawlJob{
		ID:     id,
		URL:    req.URL,
		Status: models.StatusQueued,
	}
	if err := models.DB.Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save job to database",
		})
		return
	}

	// Add job to queue
	QueueService.Enqueue(models.CrawlJob{ID: id, URL: req.URL})
	c.JSON(http.StatusOK, gin.H{
		"message": "Crawl request received",
		"id":      id,
	})
}
