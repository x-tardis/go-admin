package database

import (
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/x-tardis/go-admin/common/config"
	"github.com/x-tardis/go-admin/common/global"
)

func newMsql(source string) (*gorm.DB, error) {
	db, err := sql.Open("mysql", source)
	if err != nil {
		return nil, err
	}
	global.Cfg.SetDb(&config.DBConfig{
		Driver: "mysql",
		DB:     db,
	})
	return gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
}
