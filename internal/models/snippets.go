package models

import (
	"database/sql"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModel type wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// insert a new snippet into DB
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}

// return a specific snippet by ID
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

// return 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
