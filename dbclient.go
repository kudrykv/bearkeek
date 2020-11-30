package bearkeek

import (
	"context"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type dbclient struct {
	db *gorm.DB
}

func opendb(path string) (*gorm.DB, error) {
	dialector := sqlite.Open(path).(*sqlite.Dialector)
	dialector.DriverName = "sqlite_custom"

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: 0,
			Colorful:      true,
			LogLevel:      logger.Info,
		}),
		NamingStrategy: ns{},
	})
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	return db, nil
}

func (r dbclient) notes(ctx context.Context, q NotesQuery) ([]Note, error) {
	var notes []Note

	prep := r.db.WithContext(ctx)

	if len(q.OrderByColumns) > 0 {
		for _, column := range q.OrderByColumns {
			prep = prep.Order(clause.OrderByColumn{
				Column: clause.Column{Name: column.Name},
				Desc:   column.Desc,
			})
		}
	}

	if q.Limit > 0 {
		prep = prep.Limit(q.Limit)
	}

	if res := prep.Preload("Tags").Find(&notes); res.Error != nil {
		return nil, fmt.Errorf("notes select: %w", res.Error)
	}

	return notes, nil
}

func (r dbclient) note(ctx context.Context, id string) (Note, error) {
	var note Note

	res := r.db.WithContext(ctx).First(&note, "ZUNIQUEIDENTIFIER = ?", id)

	if res.Error != nil {
		return Note{}, fmt.Errorf("by id: %w", res.Error)
	}

	if err := r.db.Model(&note).Association("Tags").Find(&note.Tags); err != nil {
		return Note{}, fmt.Errorf("association: %w", err)
	}

	return note, nil
}
