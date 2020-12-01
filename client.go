package bearkeek

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type Bear struct {
	db dbclient
}

func NewDefault() (Bear, error) {
	path, err := compileDefaultPath()
	if err != nil {
		return Bear{}, fmt.Errorf("default path: %w", err)
	}

	return New(path)
}

func New(path string) (Bear, error) {
	db, err := opendb(path)
	if err != nil {
		return Bear{}, fmt.Errorf("open db: %w", err)
	}

	return Bear{db: dbclient{db: db}}, nil
}

type NotesQuery struct {
	Tags  []MatchingTag
	Terms []string
	Limit int
}

type MatchingTag struct {
	Name    string
	Exclude bool
}

func (b Bear) Notes(ctx context.Context, q NotesQuery) ([]Note, error) {
	notes, err := b.db.Notes(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("notes: %w", err)
	}

	return notes, nil
}

func (b Bear) Note(ctx context.Context, id string) (Note, error) {
	note, err := b.db.Note(ctx, id)
	if err != nil {
		return Note{}, fmt.Errorf("note: %w", err)
	}

	return note, nil
}

type TagsQuery struct {
	Term string
}

func (b Bear) Tags(ctx context.Context, q TagsQuery) ([]Tag, error) {
	tags, err := b.db.Tags(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("tags: %w", err)
	}

	return tags, nil
}

const defaultPath = "~/Library/Group Containers/9K33E3U3T4.net.shinyfrog.bear/Application Data/database.sqlite"

func compileDefaultPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("user home dir: %w", err)
	}

	return strings.Replace(defaultPath, "~", dir, 1), nil
}
