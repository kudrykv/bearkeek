package bearkeek

import (
	"context"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type dbclient struct {
	db *gorm.DB
}

func open(path string) (*gorm.DB, error) {
	dialector := sqlite.Open(path).(*sqlite.Dialector)
	dialector.DriverName = "sqlite_custom"

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	return db, nil
}

func (r dbclient) Notes(ctx context.Context, q NotesQuery) ([]Note, error) {
	var notes []Note
	if tx := r.db.WithContext(ctx).Find(&notes); tx.Error != nil {
		return nil, fmt.Errorf("notes select: %w", tx.Error)
	}

	return notes, nil
}

func (r dbclient) note(ctx context.Context, id string) (Note, error) {
	var note Note
	if tx := r.db.WithContext(ctx).First(&note, "ZUNIQUEIDENTIFIER = ?", id); tx.Error != nil {
		return Note{}, fmt.Errorf("by id: %w", tx.Error)
	}

	return note, nil
}
