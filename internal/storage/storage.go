package storage

import "time"

// JobStatus tracks where we are with each job posting
type JobStatus string

const (
	StatusNew     JobStatus = "new"
	StatusApplied JobStatus = "applied"
	StatusSkipped JobStatus = "skipped"
	StatusResponded JobStatus = "responded"
)

// JobRecord is what we store in the database
type JobRecord struct {
	ID            string
	Title         string
	Company       string
	Location      string
	URL           string
	Source        string
	Status        JobStatus
	FoundAt       time.Time
	AppliedAt     *time.Time
	ResumeVersion string
	Notes         string
}

// Storage is the interface for persisting job data
type Storage interface {
	// Save stores a new job, ignores duplicates
	Save(job JobRecord) error

	// Exists checks if we've seen this URL before (deduplication)
	Exists(url string) (bool, error)

	// UpdateStatus changes the status of a job
	UpdateStatus(id string, status JobStatus) error

	// List returns jobs filtered by status
	List(status JobStatus) ([]JobRecord, error)

	// Close cleans up the connection
	Close() error
}
