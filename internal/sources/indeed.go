package sources

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type IndeedSource struct {
	client *http.Client
}

func NewIndeedSource() *IndeedSource {
	return &IndeedSource{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *IndeedSource) Name() string {
	return "indeed"
}

func (s *IndeedSource) IsHealthy() bool {
	resp, err := s.client.Get("https://www.indeed.com")
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == 200
}

// rssFeed represents the Indeed RSS response
type rssItem struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
	Source  string `xml:"source"`
	Snippet string `xml:"description"`
}

type rssFeed struct {
	Items []rssItem `xml:"channel>item"`
}

func (s *IndeedSource) Search(query Query) ([]Job, error) {
	// Build Indeed RSS URL
	keywords := ""
	for i, k := range query.Keywords {
		if i > 0 {
			keywords += " "
		}
		keywords += k
	}

	rssURL := fmt.Sprintf(
		"https://www.indeed.com/rss?q=%s&l=%s&sort=date",
		url.QueryEscape(keywords),
		url.QueryEscape(query.Location),
	)

	resp, err := s.client.Get(rssURL)
	if err != nil {
		return nil, fmt.Errorf("fetch indeed rss: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("indeed rss returned status %d", resp.StatusCode)
	}

	var feed rssFeed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("decode rss: %w", err)
	}

	var jobs []Job
	for _, item := range feed.Items {
		pubDate, _ := time.Parse(time.RFC1123Z, item.PubDate)
		jobs = append(jobs, Job{
			ID:          item.Link,
			Title:       item.Title,
			URL:         item.Link,
			Description: item.Snippet,
			Source:      s.Name(),
			PostedAt:    pubDate,
			FoundAt:     time.Now(),
		})
	}

	return jobs, nil
}
