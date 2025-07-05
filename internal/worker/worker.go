package worker

import (
	"fmt"
	"time"

	"github.com/chinu-anand/crawlerx/internal/crawler"
	"github.com/chinu-anand/crawlerx/internal/models"
	"github.com/chinu-anand/crawlerx/internal/ws"
)

func Start(queue <-chan models.JobPayload) {
	for job := range queue {

		log := func(msg string) {
			fmt.Printf("%s: %s\n", job.ID, msg)
		}

		log("Starting Job")

		// Update job status to processing
		models.DB.Model(&models.CrawlJob{}).
			Where("id = ?", job.ID).
			Update("status", models.StatusProcessing)

		// Crawl the URL
		const maxRetries = 3
		const retryDelay = 3 * time.Second

		var result *crawler.CrawlResponse
		var err error

		for attempts := 0; attempts <= maxRetries; attempts++ {
			result, err = crawler.Crawl(job.URL)
			if err == nil {
				break
			}

			fmt.Printf("🔁 Retry %d for job %s: %v\n", attempts+1, job.ID, err)

			if attempts < maxRetries {
				time.Sleep(retryDelay)
			}
		}

		if err != nil {
			models.DB.Model(&models.CrawlJob{}).
				Where("id = ?", job.ID).
				Updates(models.CrawlJob{
					Status: models.StatusFailed,
					Error:  err.Error(),
				})
			ws.GetHub().BroadcastJobUpdate(job.ID, job.URL, models.StatusFailed, "")
			continue
		}

		// Update job status to done
		if err := models.DB.Model(&models.CrawlJob{}).
			Where("id = ?", job.ID).
			Updates(map[string]interface{}{
				"status":      models.StatusDone,
				"title":       result.Title,
				"description": result.Description,
				"links":       models.StringArray(result.Links),
				"error":       "",
			}).Error; err != nil {
			// If update fails, log the error
			log(fmt.Sprintf("❌ Error updating job in database: %v\n", err))

			// Try to save the error to the database
			models.DB.Model(&models.CrawlJob{}).
				Where("id = ?", job.ID).
				Updates(map[string]interface{}{
					"status": models.StatusFailed,
					"error":  fmt.Sprintf("Database error: %v", err),
				})

			ws.GetHub().BroadcastJobUpdate(job.ID, job.URL, models.StatusFailed, fmt.Sprintf("Database error: %v", err))
			log("❌ Job Failed: " + fmt.Sprintf("Database error: %v", err))
			continue
		}

		// Send result to client
		ws.GetHub().BroadcastJobUpdate(job.ID, job.URL, models.StatusDone, "")
		log("✅ Job Completed")
	}
}
