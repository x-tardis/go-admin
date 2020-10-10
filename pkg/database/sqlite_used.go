// +build sqlite3

package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func newSqlite3(source string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("sqlite3"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
}
