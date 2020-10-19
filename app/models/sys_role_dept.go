package models

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/pkg/deployed"
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

func (cRoleDept) Create(roleId int, deptIds []int) error {
	// ORM不支持批量插入所以需要拼接 sql 串
	sql := "INSERT INTO `sys_role_dept` (`role_id`,`dept_id`) VALUES "

	for i := 0; i < len(deptIds); i++ {
		if len(deptIds)-1 == i {
			//最后一条数据 以分号结尾
			sql += fmt.Sprintf("(%d,%d);", roleId, deptIds[i])
		} else {
			sql += fmt.Sprintf("(%d,%d),", roleId, deptIds[i])
		}
	}
	deployed.DB.Exec(sql)

	return nil
}

func (cRoleDept) DeleteRoleDept(roleId int) error {
	return deployed.DB.Scopes(RoleDeptDB()).
		Where("role_id=?", roleId).Delete(&RoleDept{}).Error
}
