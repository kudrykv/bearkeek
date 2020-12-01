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
	path, err := defaultPath()
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
	OrderByColumns []OrderByColumn
	Tags           []MatchingTag
	Terms          []string
	Limit          int
}

type OrderByColumn struct {
	Name string
	Desc bool
}

type MatchingTag struct {
	Name    string
	Exclude bool
}

func (b Bear) Notes(ctx context.Context, q NotesQuery) ([]Note, error) {
	notes, err := b.db.notes(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("notes: %w", err)
	}

	return notes, nil
}

func (b Bear) Note(ctx context.Context, id string) (Note, error) {
	note, err := b.db.note(ctx, id)
	if err != nil {
		return Note{}, fmt.Errorf("note: %w", err)
	}

	return note, nil
}

func defaultPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("user home dir: %w", err)
	}

	return strings.Replace(DefaultPath, "~", dir, 1), nil
}
