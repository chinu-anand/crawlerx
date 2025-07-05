package api

import (
	"encoding/csv"
	"net/http"
	"strconv"

	"github.com/chinu-anand/crawlerx/internal/models"

	"github.com/gin-gonic/gin"
)

func GetJobStatus(c *gin.Context) {
	id := c.Param("id")
	var job models.CrawlJob
	if err := models.DB.First(&job, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": job.ID, "status": job.Status, "error": job.Error})
}

func GetJobByID(c *gin.Context) {
	id := c.Param("id")
	var job models.CrawlJob
	if err := models.DB.First(&job, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	c.JSON(http.StatusOK, job)
}

func GetJobs(c *gin.Context) {
	var jobs []models.CrawlJob

	status := c.Query("status")
	limit := toInt(c.DefaultQuery("limit", "10"))
	offset := toInt(c.DefaultQuery("offset", "0"))
	search := c.Query("search")

	q := models.DB.Order("created_at desc")
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if search != "" {
		p := "%" + search + "%"
		q = q.Where("url ILIKE ? OR title ILIKE ?", p, p)
	}

	if err := q.Limit(limit).Offset(offset).Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query failed"})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

func ExportJob(c *gin.Context) {
	id := c.Param("id")
	format := c.DefaultQuery("format", "json")

	var job models.CrawlJob
	if err := models.DB.First(&job, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	switch format {
	case "json":
		c.Header("Content-Disposition", "attachment; filename=\"job_"+id+".json\"")
		c.JSON(http.StatusOK, job)

	case "csv":
		c.Header("Content-Disposition", "attachment; filename=\"job_"+id+".csv\"")
		c.Header("Content-Type", "text/csv")
		writer := csv.NewWriter(c.Writer)
		writer.Write([]string{"Field", "Value"})
		writer.Write([]string{"ID", job.ID})
		writer.Write([]string{"URL", job.URL})
		writer.Write([]string{"Title", job.Title})
		writer.Write([]string{"Description", job.Description})
		writer.Write([]string{"Status", string(job.Status)})
		for _, link := range job.Links {
			writer.Write([]string{"Link", link})
		}
		writer.Flush()

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported format"})
	}
}

func GetJobSummary(c *gin.Context) {
	type Result struct {
		Status string
		Count  int64
	}
	var results []Result
	models.DB.Model(&models.CrawlJob{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results)

	summary := map[string]int64{"total": 0}
	for _, r := range results {
		summary[r.Status] = r.Count
		summary["total"] += r.Count
	}
	c.JSON(http.StatusOK, summary)
}

func toInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}
