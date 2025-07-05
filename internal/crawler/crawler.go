package crawler

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type CrawlResponse struct {
	Title       string
	Description string
	Links       []string
}

func Crawl(url string) (*CrawlResponse, error) {
	// Make the GET Request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	// Load HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to load HTML: %w", err)
	}

	// Extract title
	title := doc.Find("title").First().Text()

	// Exract meta description
	description, _ := doc.Find(`meta[name="description"]`).Attr("content")

	// Extract links
	var links []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			links = append(links, href)
		}
	})

	return &CrawlResponse{
		Title:       title,
		Description: description,
		Links:       links,
	}, nil
}
