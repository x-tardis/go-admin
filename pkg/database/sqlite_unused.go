// +build !sqlite3

package database

import "gorm.io/gorm"

func newSqlite3(string) (*gorm.DB, error) {
	panic("please build tags with sqlite3!")
}
