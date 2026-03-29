package gormigrate

import (
	"gorm.io/gorm"
)

func TableExists(db *gorm.DB, schema, table string) bool {
	var count int64
	db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ?",
		schema, table).Scan(&count)
	return count > 0
}

func ColumnExists(db *gorm.DB, schema, table, column string) bool {
	var count int64
	db.Raw("SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = ? AND table_name = ? AND column_name = ?",
		schema, table, column).Scan(&count)
	return count > 0
}

func IndexExists(db *gorm.DB, schema, table, index string) bool {
	var count int64
	db.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = ? AND table_name = ? AND index_name = ?",
		schema, table, index).Scan(&count)
	return count > 0
}

func ConstraintExists(db *gorm.DB, schema, table, constraint string) bool {
	var count int64
	db.Raw("SELECT COUNT(*) FROM information_schema.table_constraints WHERE table_schema = ? AND table_name = ? AND constraint_name = ?",
		schema, table, constraint).Scan(&count)
	return count > 0
}
