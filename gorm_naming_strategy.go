package bearkeek

import "gorm.io/gorm/schema"

type ns struct{}

func (ns) TableName(table string) string {
	return table
}

func (ns) ColumnName(table, column string) string {
	return table + column
}

func (ns) JoinTableName(joinTable string) string {
	return joinTable
}

func (ns) RelationshipFKName(relationship schema.Relationship) string {
	return relationship.Name
}

func (ns) CheckerName(table, column string) string {
	return table + column
}

func (ns) IndexName(table, column string) string {
	return table + column
}
