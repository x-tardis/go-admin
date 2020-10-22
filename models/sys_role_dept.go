package models

import (
	"context"

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

func RoleDeptDB() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(RoleDept{})
	}
}

type cRoleDept struct{}

var CRoleDept = new(cRoleDept)

func (cRoleDept) Create(_ context.Context, roleId int, deptIds []int) error {
	newItems := make([]RoleDept, 0, len(deptIds))
	for _, v := range deptIds {
		newItems = append(newItems, RoleDept{roleId, v})
	}
	return dao.DB.Scopes(RoleDeptDB()).Create(&newItems).Error
}

func (cRoleDept) Delete(_ context.Context, roleId int) error {
	return dao.DB.Scopes(RoleDeptDB()).
		Where("role_id=?", roleId).Delete(&RoleDept{}).Error
}
