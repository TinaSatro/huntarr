# Huntarr

## The Problem

Job searching is a full-time job. Every day: open LinkedIn, Indeed, Glassdoor, Jobright — check for new postings, evaluate fit, rewrite resume headline and summary, craft a recruiter message. Repeat. **2+ hours daily of mostly mechanical work.**

The frustrating part: the creative work (am I a fit? what's my angle for this role?) takes 10 minutes. The mechanical work (find, filter, reformat, track) takes the rest.

## What Huntarr Does

Huntarr automates the mechanical parts so you can focus on the human parts.

- **Finds** relevant jobs across multiple boards (LinkedIn, Indeed, Glassdoor, Jobright, and more)
- **Deduplicates** — never shows you the same posting twice, tracks what you've seen, applied to, or skipped
- **Adapts** your resume headline and summary to match each job description via Claude API
- **Drafts** a short recruiter outreach message per posting
- **Tracks** everything — company, role, date, status, which resume version you sent

You still apply. You still decide. Huntarr removes the noise.

## Architecture

The core design principle: **job boards change constantly**. URLs break, layouts shift, new platforms emerge. The system is built to absorb that change in one place.
cmd/            → CLI entry point and flags
internal/
sources/      → JobSource interface + implementations per board
resume/       → Resume adaptation via Claude API
outreach/     → Recruiter message generation
sender/       → Sender interface (email / browser handoff / export)
storage/      → SQLite-backed job tracking and deduplication
config/       → YAML configuration, see config.example.yaml

## Design Decisions

**Why Go**
Compiled binary, easy to distribute, strong concurrency for parallel board scraping. No runtime dependencies — runs anywhere.

**Why a JobSource interface**
Every job board is a different implementation of the same contract: `Search(query Query) ([]Job, error)`. Adding a new board means adding one file. Nothing else changes. When a board breaks — and they will — you fix one file.

**Why SQLite, not Postgres**
No server to run. Single file. Portable. If this ever becomes a service, swapping the storage layer is the only change needed — the interface stays the same.

**Why Claude API for resume adaptation**
Resume rewriting is context-sensitive and hard to template. Prompting a model with your base resume + job description produces better results than any string substitution approach. Cost per adaptation: fractions of a cent.

## Status

🚧 Active development

## Configuration

Copy `config.example.yaml` and fill in your values:
```bash
cp config/config.example.yaml config/config.yaml
```

`config.yaml` is gitignored — your API keys and resume stay local.
