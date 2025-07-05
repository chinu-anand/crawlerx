package api

type CrawlRequest struct {
	URL string `json:"url" binding:"required,url"`
}

type JobPayload struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}
