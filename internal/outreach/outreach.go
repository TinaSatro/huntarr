package outreach

// Message is a recruiter outreach message generated for a specific job
type Message struct {
	Subject string
	Body    string
	JobURL  string
	Company string
}

// Generator creates recruiter outreach messages
type Generator interface {
	// Generate creates a short recruiter message based on job details
	Generate(jobTitle string, company string, jobDescription string) (Message, error)
}
