package models

import (
	"context"

	"github.com/x-tardis/go-admin/pkg/deployed"
)

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Code     string `form:"code" json:"code" binding:"required"`
	UUID     string `form:"uuid" json:"uuid" binding:"required"`
}

func (u *Login) GetUser() (user SysUser, role SysRole, err error) {
	user, err = new(CallUser).GetWithName(context.Background(), u.Username)
	if err != nil {
		return
	}
	// check password
	err = deployed.Verify.Compare(u.Password, "", user.Password)
	if err != nil {
		return
	}
	role, err = new(CallRole).Get(context.Background(), user.RoleId)
	return
}
