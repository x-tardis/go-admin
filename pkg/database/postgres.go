package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newPostgres(source string) gorm.Dialector {
	return postgres.New(postgres.Config{
		DSN: source,
	})
}
