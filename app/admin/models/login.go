package models

import (
	"github.com/x-tardis/go-admin/pkg/deployed"
)

type Login struct {
	Username string `form:"UserName" json:"username" binding:"required"`
	Password string `form:"Password" json:"password" binding:"required"`
	Code     string `form:"Code" json:"code" binding:"required"`
	UUID     string `form:"UUID" json:"uuid" binding:"required"`
}

func (u *Login) GetUser() (user SysUser, role SysRole, e error) {
	e = deployed.DB.Table("sys_user").Where("username = ? ", u.Username).Find(&user).Error
	if e != nil {
		return
	}

	e = deployed.Verify.Compare(u.Password, "", user.Password)
	if e != nil {
		return
	}
	e = deployed.DB.Table("sys_role").Where("role_id = ? ", user.RoleId).First(&role).Error
	if e != nil {
		return
	}
	return
}
