package migrate

import (
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/models"
	"github.com/x-tardis/go-admin/models/tools"
)

func BaseModelsTables(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		err := tx.AutoMigrate(
			new(models.CasbinRule),
			new(models.User),
			new(models.Menu),
			new(models.Role),
			new(models.RoleMenu),
			new(models.RoleDept),
			new(models.Dept),
			new(models.Post),
			new(models.DictData),
			new(models.DictType),
			new(models.Config),
			new(models.Setting),
			new(models.Job),
			new(models.LoginLog),
			new(models.OperLog),
			new(models.FileDir),
			new(models.FileInfo),
			new(models.Category),
			new(models.Content),
			new(tools.SysTables),
			new(tools.SysColumns),
		)
		if err != nil {
			return err
		}
		return tx.Create(&models.Migration{Version: version}).Error
	})
}
