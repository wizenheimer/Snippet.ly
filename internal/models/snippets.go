package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   sql.NullString // string
	Content sql.NullString // string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

// Insert Snippet into the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP, INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get a snippet by ID
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires >	UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// Get 10 latest snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10
	`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}
	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
