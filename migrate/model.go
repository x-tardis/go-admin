package migrate

import (
	"github.com/x-tardis/go-admin/models"
	"github.com/x-tardis/go-admin/models/tools"
)

// baseModels 所有基础模型
func baseModels() []interface{} {
	return []interface{}{
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
	}
}
