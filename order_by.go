package bearkeek

// nolint: gochecknoglobals
//goland:noinspection GoUnusedGlobalVariable
var (
	CreatedAtAsc  = OrderByColumn{Name: "ZCREATIONDATE", Desc: false}
	CreatedAtDesc = OrderByColumn{Name: "ZCREATIONDATE", Desc: true}
	UpdatedAtAsc  = OrderByColumn{Name: "ZMODIFICATIONDATE", Desc: false}
	UpdatedAtDesc = OrderByColumn{Name: "ZMODIFICATIONDATE", Desc: true}
)

func OrderByColumns(fields ...OrderByColumn) []OrderByColumn {
	return fields
}
