package models

import (
	"context"

	"github.com/thinkgos/x/lib/ternary"

	"github.com/x-tardis/go-admin/deployed"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

// Login 登录
type Login struct {
	Username string `form:"username" json:"username" binding:"required"` // 用户名
	Password string `form:"password" json:"password" binding:"required"` // 密码
	CID      string `form:"cid" json:"cid" binding:"required"`           // 验证码id
	CCode    string `form:"code" json:"ccode" binding:"required"`        // 验证码
}

// Get 获取 user 和 role 和用户是否已使能
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
		DeptId:    user.DeptId,
		PostId:    user.PostId,
		RoleId:    role.RoleId,
		RoleName:  role.RoleName,
		RoleKey:   ternary.IfString(user.Username == SuperRoot, SuperRoot, role.RoleKey),
		DataScope: role.DataScope,
	}, user.Status == StatusEnable, nil
}
