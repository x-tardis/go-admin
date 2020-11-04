package models

import (
	"context"

	"github.com/thinkgos/sharp/iorm/trans"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
)

// RoleDept role dept关系表
type RoleDept struct {
	RoleId int
	DeptId int
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
		Delete(&RoleDept{}, "role_id=?", roleId).Error
}

// DeleteWithDept 通过部门id删除
func (cRoleDept) DeleteWithDept(ctx context.Context, deptId int) error {
	return dao.DB.Scopes(RoleDeptDB(ctx)).
		Delete(&RoleDept{}, "dept_id=?", deptId).Error
}

// GetDeptTreeOption 获取部门树和角色已选的部门id列表
func (cRoleDept) GetDeptTreeOption(ctx context.Context, roleId int) ([]DeptNameLabel, []int, error) {
	tree, err := CDept.QueryLabelTree(ctx)
	if err != nil {
		return nil, nil, err
	}
	deptIds := make([]int, 0)
	if roleId != 0 {
		deptIds, err = CRole.GetDeptIds(ctx, roleId)
	}
	return tree, deptIds, err
}
