package bearkeek

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type dbclient struct {
	db *gorm.DB
}

func opendb(path string) (*gorm.DB, error) {
	dialector := sqlite.Open(path).(*sqlite.Dialector)
	dialector.DriverName = "sqlite_custom"

	db, err := gorm.Open(dialector, &gorm.Config{NamingStrategy: ns{}})
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	return db, nil
}

func (r dbclient) Notes(ctx context.Context, q NotesQuery) ([]Note, error) {
	var notes []Note

	prep := r.db.WithContext(ctx).Preload("Tags")

	var err error
	if prep, err = notesTags(r.db, prep, q); err != nil {
		return nil, fmt.Errorf("notes tags: %w", err)
	}

	prep = r.notesWhereTerms(prep, q)
	prep = r.notesLimit(prep, q)

	res := prep.
		Where("ZARCHIVED = 0").Where("ZTRASHED = 0").
		Order("ZCREATIONDATE desc").
		Find(&notes)

	if res.Error != nil {
		return nil, fmt.Errorf("notes select: %w", res.Error)
	}

	return notes, nil
}

func (r dbclient) notesWhereTerms(prep *gorm.DB, q NotesQuery) *gorm.DB {
	if len(q.Terms) == 0 {
		prep.Select("ZSFNOTE.*, 0 as titlehit")

		return prep
	}

	where := r.db.Where("utf8lower(ZTEXT) like utf8lower(?)", "%"+q.Terms[0]+"%")
	slect := []string{"ZSFNOTE.*, utf8lower(ZTITLE) like utf8lower(?)"}
	params := []interface{}{"%" + q.Terms[0] + "%"}

	for i := 1; i < len(q.Terms); i++ {
		where = where.Or("utf8lower(ZTEXT) like utf8lower(?)", "%"+q.Terms[i]+"%")

		slect = append(slect, "utf8lower(ZTITLE) like utf8lower(?)")
		params = append(params, "%"+q.Terms[i]+"%")
	}

	prep = prep.Select(strings.Join(slect, " OR ")+" as titlehit", params...).
		Order("titlehit desc")

	return prep.Where(where)
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
		if tl.exact(exclude, true) && tl.exact(include, false) {
			noteIDs = append(noteIDs, nn.ID)
		}
	}

	return prep.Where("Z_PK IN (?)", noteIDs), nil
}

func (r dbclient) notesLimit(prep *gorm.DB, q NotesQuery) *gorm.DB {
	if q.Limit <= 0 {
		return prep
	}

	return prep.Limit(q.Limit)
}

func (r dbclient) Note(ctx context.Context, id string) (Note, error) {
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

func (r dbclient) Tags(ctx context.Context, q TagsQuery) ([]Tag, error) {
	var tags []Tag

	res := r.db.WithContext(ctx).Table("ZSFNOTETAG").
		Distinct("ZSFNOTETAG.Z_PK, ZSFNOTETAG.ZTITLE").
		Joins("join Z_7TAGS on Z_7TAGS.Z_14TAGS = ZSFNOTETAG.Z_PK").
		Joins("join ZSFNOTE on ZSFNOTE.Z_PK = Z_7TAGS.Z_7NOTES").
		Where("ZSFNOTE.ZTRASHED = 0 and ZSFNOTE.ZARCHIVED = 0").
		Where("ZSFNOTETAG.ZTITLE like ?", "%"+q.Term+"%").
		Find(&tags)

	if res.Error != nil {
		return nil, fmt.Errorf("find: %w", res.Error)
	}

	return tags, nil
}
