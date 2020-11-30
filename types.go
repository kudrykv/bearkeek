package bearkeek

type Note struct {
	ID                 int64  `gorm:"column:Z_PK;primaryKey"`
	UUID               string `gorm:"column:ZUNIQUEIDENTIFIER"`
	Title              string `gorm:"column:ZTITLE"`
	Subtitle           string `gorm:"column:ZSUBTITLE"`
	Archived           int    `gorm:"column:ZARCHIVED"`
	Trashed            int    `gorm:"column:ZTRASHED"`
	PermanentlyDeleted int    `gorm:"column:ZPERMANENTLYDELETED"`

	Tags []Tag `gorm:"many2many:Z_7TAGS;joinForeignKey:Z_7NOTES;joinReferences:Z_14TAGS"`
}

func (Note) TableName() string {
	return "ZSFNOTE"
}

type Tag struct {
	ID   int64  `gorm:"column:Z_PK;primaryKey"`
	Name string `gorm:"column:ZTITLE"`

	Notes []Note `gorm:"many2many:Z_7TAGS;joinForeignKey:Z_14TAGS;joinReferences:Z_7NOTES"`
}

func (Tag) TableName() string {
	return "ZSFNOTETAG"
}
