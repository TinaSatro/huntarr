package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/TinaSatro/huntarr/internal/outreach"
	"github.com/TinaSatro/huntarr/internal/resume"
	"github.com/TinaSatro/huntarr/internal/storage"
)

func main() {
	// CLI flags
	dbPath := flag.String("db", "huntarr.db", "Path to SQLite database")
	claudeKey := flag.String("claude-key", os.Getenv("CLAUDE_API_KEY"), "Claude API key")
	flag.Parse()

	if *claudeKey == "" {
		log.Fatal("Claude API key required: set CLAUDE_API_KEY or use -claude-key flag")
	}

	// Initialize storage
	store, err := storage.NewSQLiteStorage(*dbPath)
	if err != nil {
		log.Fatalf("storage init: %v", err)
	}
	defer store.Close()

	// Initialize Claude adapters
	resumeAdapter := resume.NewClaudeAdapter(*claudeKey)
	outreachGenerator := outreach.NewClaudeGenerator(*claudeKey)

	fmt.Println("Huntarr initialized successfully")
	fmt.Printf("Database: %s\n", *dbPath)

	// Smoke test — adapt a sample resume
	base := resume.BaseResume{
		Headline: "Senior SRE | Kubernetes | Infrastructure",
		Summary:  "Infrastructure architect with expertise in Kubernetes, RKE2, and airgap edge deployments.",
	}

	adapted, err := resumeAdapter.Adapt(base,
		"Staff Infrastructure Architect",
		"Acme Corp",
		"Looking for a staff-level architect with Kubernetes and edge computing experience.")
	if err != nil {
		log.Fatalf("resume adapt: %v", err)
	}

	fmt.Println("\n--- Adapted Resume ---")
	fmt.Printf("Headline: %s\n", adapted.Headline)
	fmt.Printf("Summary:  %s\n", adapted.Summary)

	// Smoke test — generate outreach message
	msg, err := outreachGenerator.Generate(
		"Staff Infrastructure Architect",
		"Acme Corp",
		"Looking for a staff-level architect with Kubernetes and edge computing experience.")
	if err != nil {
		log.Fatalf("outreach generate: %v", err)
	}

	fmt.Println("\n--- Recruiter Message ---")
	fmt.Printf("Subject: %s\n", msg.Subject)
	fmt.Printf("Body:\n%s\n", msg.Body)

	_ = store
}
