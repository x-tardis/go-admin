package models

import (
	"context"

	"github.com/thinkgos/sharp/iorm/trans"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
)

//sys_role_dept
type RoleDept struct {
	RoleId int `gorm:""`
	DeptId int `gorm:""`
}

func (RoleDept) TableName() string {
	return "sys_role_dept"
}

func RoleDeptDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(RoleDept{})
	}
}

type cRoleDept struct{}

var CRoleDept = cRoleDept{}

func (cRoleDept) Create(ctx context.Context, roleId int, deptIds []int) error {
	newItems := make([]RoleDept, 0, len(deptIds))
	for _, v := range deptIds {
		newItems = append(newItems, RoleDept{roleId, v})
	}
	return dao.DB.Scopes(RoleDeptDB(ctx)).Create(&newItems).Error
}

func (cRoleDept) DeleteWithRole(ctx context.Context, roleId int) error {
	return dao.DB.Scopes(RoleDeptDB(ctx)).
		Where("role_id=?", roleId).Delete(&RoleDept{}).Error
}
