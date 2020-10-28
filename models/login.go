package models

import (
	"context"

	"github.com/x-tardis/go-admin/deployed"
)

// Login 登录
type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	CID      string `form:"cid" json:"cid" binding:"required"`    // 验证码id
	CCode    string `form:"code" json:"ccode" binding:"required"` // 验证码
}

// Get 获取user 和 role
func (sf *Login) Get() (user User, role Role, err error) {
	user, err = CUser.GetWithName(context.Background(), sf.Username)
	if err != nil {
		return
	}
	err = deployed.Verify.Compare(sf.Password, user.Salt, user.Password)
	if err != nil {
		return
	}
	role, err = CRole.Get(context.Background(), user.RoleId)
	return
}
