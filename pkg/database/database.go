package database

import (
	"gorm.io/gorm"
)

// Database interface
type Database interface {
	Setup()
	Open(conn string, cfg *gorm.Config) (db *gorm.DB, err error)
	GetConnect() string
	GetDriver() string
}

func Setup(driver string) {
	switch driver {
	case "mysql":
		var db = new(Mysql)
		db.Setup()
	case "postgres":
		var db = new(PgSql)
		db.Setup()

	case "sqlite3":
		var db = new(SqLite)
		db.Setup()
	default:
		panic("please select database driver one of [mysql|postgres|mysqlite3], if use sqlite3,build tags with sqlite3!")
	}
}
