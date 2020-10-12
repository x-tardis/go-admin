package system

import (
	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/admin/models"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/tools"
)

func GetInfo(c *gin.Context) {
	var roles = make([]string, 1)
	roles[0] = jwtauth.RoleKey(c)

	var permissions = make([]string, 1)
	permissions[0] = "*:*:*"

	var buttons = make([]string, 1)
	buttons[0] = "*:*:*"

	RoleMenu := models.RoleMenu{}
	RoleMenu.RoleId = jwtauth.RoleId(c)

	var mp = make(map[string]interface{})
	mp["roles"] = roles
	if jwtauth.RoleKey(c) == "admin" || jwtauth.RoleKey(c) == "系统管理员" {
		mp["permissions"] = permissions
		mp["buttons"] = buttons
	} else {
		list, _ := RoleMenu.GetPermis()
		mp["permissions"] = list
		mp["buttons"] = list
	}

	sysuser := models.SysUser{}
	sysuser.UserId = jwtauth.UserId(c)
	user, err := sysuser.Get()
	tools.HasError(err, "", 500)

	mp["introduction"] = " am a super administrator"

	mp["avatar"] = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	if user.Avatar != "" {
		mp["avatar"] = user.Avatar
	}
	mp["userName"] = user.NickName
	mp["userId"] = user.UserId
	mp["deptId"] = user.DeptId
	mp["name"] = user.NickName

	servers.OKWithRequestID(c, mp, "")
}
