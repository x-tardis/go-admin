// +build !sqlite3

package database

import (
	"errors"

	"gorm.io/gorm"
)

type SqLite struct{}

func (e *SqLite) Setup() {
	panic("please build tags with sqlite3!")
}

// 打开数据库连接
func (*SqLite) Open(conn string, cfg *gorm.Config) (db *gorm.DB, err error) {
	return nil, errors.New("please build tags with sqlite3!")
}

func (e *SqLite) GetConnect() string { return "" }

func (e *SqLite) GetDriver() string { return "" }
