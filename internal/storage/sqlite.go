package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type SQLiteStorage struct {
	db *sql.DB
}

// NewSQLiteStorage opens the database and creates tables if needed
func NewSQLiteStorage(path string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("create tables: %w", err)
	}

	return &SQLiteStorage{db: db}, nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS jobs (
			id             TEXT PRIMARY KEY,
			title          TEXT NOT NULL,
			company        TEXT NOT NULL,
			location       TEXT,
			url            TEXT UNIQUE NOT NULL,
			source         TEXT,
			status         TEXT DEFAULT 'new',
			found_at       DATETIME,
			applied_at     DATETIME,
			resume_version TEXT,
			notes          TEXT
		)
	`)
	return err
}

func (s *SQLiteStorage) Save(job JobRecord) error {
	_, err := s.db.Exec(`
		INSERT OR IGNORE INTO jobs
			(id, title, company, location, url, source, status, found_at)
		VALUES
			(?, ?, ?, ?, ?, ?, ?, ?)
	`, job.ID, job.Title, job.Company, job.Location,
		job.URL, job.Source, StatusNew, time.Now())
	return err
}

func (s *SQLiteStorage) Exists(url string) (bool, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM jobs WHERE url = ?`, url).Scan(&count)
	return count > 0, err
}

func (s *SQLiteStorage) UpdateStatus(id string, status JobStatus) error {
	_, err := s.db.Exec(`UPDATE jobs SET status = ? WHERE id = ?`, status, id)
	return err
}

func (s *SQLiteStorage) List(status JobStatus) ([]JobRecord, error) {
	rows, err := s.db.Query(`SELECT id, title, company, location, url, source, status, found_at FROM jobs WHERE status = ?`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []JobRecord
	for rows.Next() {
		var j JobRecord
		err := rows.Scan(&j.ID, &j.Title, &j.Company, &j.Location, &j.URL, &j.Source, &j.Status, &j.FoundAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
