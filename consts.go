package bearkeek

const (
	DefaultPath = "~/Library/Group Containers/9K33E3U3T4.net.shinyfrog.bear/Application Data/database.sqlite"

	QueryAllNotes = `
SELECT
	ZUNIQUEIDENTIFIER, ZTITLE, ZSUBTITLE, ZARCHIVED, ZTRASHED, ZPERMANENTLYDELETED
FROM
	ZSFNOTE
`
)
