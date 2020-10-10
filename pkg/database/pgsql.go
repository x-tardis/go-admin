package database

import (
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/x-tardis/go-admin/common/config"
	"github.com/x-tardis/go-admin/common/global"
)

func newPgSql(source string) (*gorm.DB, error) {
	db, err := sql.Open("postgresql", source)
	if err != nil {
		return nil, err
	}
	global.Cfg.SetDb(&config.DBConfig{
		Driver: "mysql",
		DB:     db,
	})
	return gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
}
