package bearkeek

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
)

type Client struct {
	db *sql.DB
}

func New() *Client {
	return &Client{}
}

func (r *Client) Open(path string) error {
	var err error
	r.db, err = sql.Open("sqlite_custom", path)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	return nil
}

func (r *Client) OpenDefault() error {
	dir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("user home dir: %w", err)
	}

	replace := strings.Replace(DefaultPath, "~", dir, 1)
	if err := r.Open(replace); err != nil {
		return fmt.Errorf("open default: %w", err)
	}

	return nil
}

func (r *Client) AllNotes(ctx context.Context) ([]Note, error) {
	qc, err := r.db.QueryContext(ctx, QueryAllNotes)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}

	notes := make([]Note, 0, 100)

	for qc.Next() {
		if err = qc.Err(); err != nil {
			return nil, fmt.Errorf("next: %w", err)
		}

		var n Note
		if err = qc.Scan(&n.ID, &n.Title, &n.Subtitle, &n.Archived, &n.Trashed, &n.PermanentlyDeleted); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		notes = append(notes, n)
	}

	return notes, nil
}

func (r *Client) GetNoteByID(ctx context.Context, id string) (Note, error) {
	row := r.db.QueryRowContext(ctx, QueryANote, id)

	var n Note
	if err := row.Scan(&n.ID, &n.Title, &n.Subtitle, &n.Archived, &n.Trashed, &n.PermanentlyDeleted); err != nil {
		return Note{}, fmt.Errorf("scan: %w", err)
	}

	if err := row.Err(); err != nil {
		return Note{}, fmt.Errorf("row err: %w", err)
	}

	return n, nil
}
