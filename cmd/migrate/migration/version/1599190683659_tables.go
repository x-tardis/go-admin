package version

import (
	"runtime"

	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/app/models/tools"
	"github.com/x-tardis/go-admin/cmd/migrate/migration"
	common "github.com/x-tardis/go-admin/common/models"
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
		new(models.SysJob),
		new(models.Config),
		new(models.Setting),
		new(models.SysFileDir),
		new(models.SysFileInfo),
		new(models.Category),
		new(models.SysContent),
	)
	if err != nil {
		return err
	}
	return db.Create(&common.Migration{
		Version: version,
	}).Error
}
