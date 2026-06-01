package resume

// BaseResume holds your original resume content
type BaseResume struct {
	Headline string
	Summary  string
	FullText string
}

// AdaptedResume is the result of tailoring for a specific job
type AdaptedResume struct {
	Headline string
	Summary  string
	JobURL   string
	JobTitle string
	Company  string
}

// Adapter takes your base resume and adapts it for a specific job
type Adapter interface {
	// Adapt rewrites headline and summary to match the job description
	Adapt(base BaseResume, jobTitle string, company string, jobDescription string) (AdaptedResume, error)
}
