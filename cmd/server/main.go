package main

import (
	"fmt"
	"log"

	"github.com/chinu-anand/crawlerx/configs"
	"github.com/chinu-anand/crawlerx/internal/api"
	"github.com/chinu-anand/crawlerx/internal/models"
	"github.com/chinu-anand/crawlerx/internal/queue"
	"github.com/chinu-anand/crawlerx/internal/worker"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	configs.LoadConfig()

	// Initialize database
	models.InitDB()

	// Create Job Queue
	jobQueue := queue.NewRedisQueue(configs.App.RedisURL)
	api.QueueService = jobQueue

	// Start worker
	go worker.Start(jobQueue.Dequeue())

	// Setup router
	r := gin.Default()
	api.RegisterRoutes(r.Group("/api/v1"))

	// start the server
	port := configs.App.Port

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

	fmt.Printf("Server is running on port %s\n", port)
}
