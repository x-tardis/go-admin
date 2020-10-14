package models

import (
	"fmt"

	"github.com/x-tardis/go-admin/pkg/deployed"
)

//sys_role_dept
type SysRoleDept struct {
	RoleId int `gorm:""`
	DeptId int `gorm:""`
}

func (SysRoleDept) TableName() string {
	return "sys_role_dept"
}

func (rm *SysRoleDept) Insert(roleId int, deptIds []int) (bool, error) {
	//ORM不支持批量插入所以需要拼接 sql 串
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

	return true, nil
}

func (rm *SysRoleDept) DeleteRoleDept(roleId int) (bool, error) {
	if err := deployed.DB.Table("sys_role_dept").Where("role_id = ?", roleId).Delete(&rm).Error; err != nil {
		return false, err
	}
	var role SysRole
	if err := deployed.DB.Table("sys_role").Where("role_id = ?", roleId).First(&role).Error; err != nil {
		return false, err
	}

	return true, nil

}
