package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	r.GET("/ws", WebSocketHandler)

	r.POST("/crawl", PostCrawl)

	r.GET("/status/:id", GetJobStatus)

	r.GET("/jobs", GetJobs)

	r.GET("/jobs/:id", GetJobByID)

	r.GET("/jobs/export/:id", ExportJob)

	r.GET("/jobs/summary", GetJobSummary)
}
