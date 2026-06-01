package sources

import "time"

// Job represents a job posting found on any board
type Job struct {
	ID          string
	Title       string
	Company     string
	Location    string
	URL         string
	Description string
	Source      string
	PostedAt    time.Time
	FoundAt     time.Time
}

// Query represents search parameters
type Query struct {
	Keywords        []string
	Location        string
	ExperienceLevel string
}

// JobSource is the interface every job board must implement.
// Adding a new board = adding one file that implements this interface.
// Nothing else in the system changes.
type JobSource interface {
	// Name returns the board identifier e.g. "linkedin", "indeed"
	Name() string

	// Search finds jobs matching the query
	Search(query Query) ([]Job, error)

	// IsHealthy checks if the source is reachable and working
	IsHealthy() bool
}
