package bearkeek

type Note struct {
	ID                 string `gorm:"column:ZUNIQUEIDENTIFIER"`
	Title              string `gorm:"column:ZTITLE"`
	Subtitle           string `gorm:"column:ZSUBTITLE"`
	Archived           int    `gorm:"column:ZARCHIVED"`
	Trashed            int    `gorm:"column:ZTRASHED"`
	PermanentlyDeleted int    `gorm:"column:ZPERMANENTLYDELETED"`
}

func (Note) TableName() string {
	return "ZSFNOTE"
}
