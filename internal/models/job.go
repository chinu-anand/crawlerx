package models

import (
	"time"
)

type JobStatus string

const (
	StatusQueued     JobStatus = "queued"
	StatusProcessing JobStatus = "processing"
	StatusDone       JobStatus = "done"
	StatusFailed     JobStatus = "failed"
)

type CrawlJob struct {
	ID          string `gorm:"primaryKey"`
	URL         string
	Title       string
	Description string
	Links       StringArray `gorm:"type:json"`
	Status      JobStatus
	Error       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type JobPayload struct {
	ID      string `json:"id"`
	URL     string `json:"url"`
	Retries int
}
