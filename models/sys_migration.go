package models

import (
	"context"
	"time"

	"github.com/thinkgos/sharp/iorm/trans"
	"gorm.io/gorm"
)

type Migration struct {
	Version   string    `gorm:"primary_key"`
	ApplyTime time.Time `gorm:"autoCreateTime"`
}

// TableName implement schema.Tabler interface
func (Migration) TableName() string {
	return "sys_migration"
}

// MigrationDB migration db scopes
func MigrationDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(Migration{})
	}
}

func GetMigrationCount(ctx context.Context, db *gorm.DB, ver int) (count int64, err error) {
	err = db.Scopes(MigrationDB(ctx)).
		Where("version=?", ver).Count(&count).Error
	return
}
