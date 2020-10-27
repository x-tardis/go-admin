package models

import (
	"context"

	"github.com/thinkgos/sharp/iorm/trans"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
)

// RoleDept role dept关系表
type RoleDept struct {
	RoleId int `gorm:""`
	DeptId int `gorm:""`
}

// TableName implement schema.Tabler interface
func (RoleDept) TableName() string {
	return "sys_role_dept"
}

// RoleDeptDB role dept db scopes
func RoleDeptDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(RoleDept{})
	}
}

type cRoleDept struct{}

// CRoleDept 实例
var CRoleDept = cRoleDept{}

// BatchCreate 批量创建
func (cRoleDept) BatchCreate(ctx context.Context, roleId int, deptIds []int) error {
	newItems := make([]RoleDept, 0, len(deptIds))
	for _, deptId := range deptIds {
		newItems = append(newItems, RoleDept{roleId, deptId})
	}
	return dao.DB.Scopes(RoleDeptDB(ctx)).Create(&newItems).Error
}

// DeleteWithRole 通过角色id删除
func (cRoleDept) DeleteWithRole(ctx context.Context, roleId int) error {
	return dao.DB.Scopes(RoleDeptDB(ctx)).
		Where("role_id=?", roleId).Delete(&RoleDept{}).Error
}
