package bearkeek

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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

	prep := r.db.WithContext(ctx).Preload("Tags")

	prep = notesOrder(prep, q)
	prep = notesLimit(prep, q)

	var err error
	if prep, err = notesTags(r.db, prep, q); err != nil {
		return nil, fmt.Errorf("notes tags: %w", err)
	}

	if res := prep.Find(&notes); res.Error != nil {
		return nil, fmt.Errorf("notes select: %w", res.Error)
	}

	return notes, nil
}

func notesTags(db *gorm.DB, prep *gorm.DB, q NotesQuery) (*gorm.DB, error) {
	if len(q.Tags) == 0 {
		return prep, nil
	}

	var nat []struct {
		ID   int64  `gorm:"column:id"`
		Tags string `gorm:"column:tags"`
	}

	res := db.Table("ZSFNOTETAG").
		Joins("JOIN `Z_7TAGS` ON `Z_7TAGS`.`Z_14TAGS` = `ZSFNOTETAG`.`Z_PK`").
		Select("Z_7NOTES as id, group_concat(utf8lower(ZTITLE), '###') as tags").
		Group("Z_7NOTES").Find(&nat)
	if res.Error != nil {
		return nil, fmt.Errorf("tags for notes: %w", res.Error)
	}

	include := make([]string, 0, len(q.Tags))
	exclude := make([]string, 0, len(q.Tags))

	for _, tag := range q.Tags {
		if tag.Exclude {
			exclude = append(exclude, tag.Name)
		} else {
			include = append(include, tag.Name)
		}
	}

	var noteIDs []int64

	for _, nn := range nat {
		tl := taglist(strings.Split(nn.Tags, "###"))
		if tl.excludeAll(exclude) && tl.includeAll(include) {
			noteIDs = append(noteIDs, nn.ID)
		}
	}

	return prep.Where("Z_PK IN (?)", noteIDs), nil
}

type taglist []string

func (t taglist) includeAll(in []string) bool {
	for _, left := range in {
		hit := false

		for _, right := range t {
			if left == right {
				hit = true

				break
			}
		}

		if !hit {
			return false
		}
	}

	return true
}

func (t taglist) excludeAll(in []string) bool {
	for _, left := range in {
		hit := false

		for _, right := range t {
			if left == right {
				hit = true

				break
			}
		}

		if hit {
			return false
		}
	}

	return true
}

func notesOrder(prep *gorm.DB, q NotesQuery) *gorm.DB {
	if len(q.OrderByColumns) == 0 {
		return prep
	}

	for _, column := range q.OrderByColumns {
		prep = prep.Order(clause.OrderByColumn{
			Column: clause.Column{Name: column.Name},
			Desc:   column.Desc,
		})
	}

	return prep
}

func notesLimit(prep *gorm.DB, q NotesQuery) *gorm.DB {
	if q.Limit <= 0 {
		return prep
	}

	return prep.Limit(q.Limit)
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
