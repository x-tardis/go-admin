package models

import (
	"context"

	"github.com/x-tardis/go-admin/deployed"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

// Login 登录
type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	CID      string `form:"cid" json:"cid" binding:"required"`    // 验证码id
	CCode    string `form:"code" json:"ccode" binding:"required"` // 验证码
}

// Get 获取user 和 role
func (sf *Login) Get() (jwtauth.Identities, bool, error) {
	user, err := CUser.GetWithName(context.Background(), sf.Username)
	if err != nil {
		return jwtauth.Identities{}, false, err
	}
	err = deployed.Verify.Compare(sf.Password, user.Salt, user.Password)
	if err != nil {
		return jwtauth.Identities{}, false, err
	}
	role, err := CRole.Get(context.Background(), user.RoleId)
	if err != nil {
		return jwtauth.Identities{}, false, err
	}
	return jwtauth.Identities{
		UserId:    user.UserId,
		Username:  user.Username,
		RoleId:    role.RoleId,
		RoleName:  role.RoleName,
		RoleKey:   role.RoleKey,
		DataScope: role.DataScope,
	}, user.Status == StatusEnable, nil
}
