package version

import (
	"runtime"

	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/app/models/tools"
	"github.com/x-tardis/go-admin/cmd/migrate/migration"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1599190683659Tables)
}

func _1599190683659Tables(db *gorm.DB, version string) error {
	err := db.Debug().Migrator().AutoMigrate(
		new(models.CasbinRule),
		new(models.Dept),
		new(models.Config),
		new(tools.SysTables),
		new(tools.SysColumns),
		new(models.Menu),
		new(models.LoginLog),
		new(models.OperLog),
		new(models.RoleMenu),
		new(models.RoleDept),
		new(models.User),
		new(models.Role),
		new(models.Post),
		new(models.DictData),
		new(models.DictType),
		new(models.Job),
		new(models.Config),
		new(models.Setting),
		new(models.FileDir),
		new(models.FileInfo),
		new(models.Category),
		new(models.Content),
	)
	if err != nil {
		return err
	}
	return db.Create(&models.Migration{
		Version: version,
	}).Error
}
