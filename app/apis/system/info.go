package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/x-tardis/go-admin/models"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/servers"
)

func GetInfo(c *gin.Context) {
	var roles = make([]string, 1)
	rolekey := jwtauth.FromRoleKey(gcontext.Context(c))
	roles[0] = rolekey

	var permissions = make([]string, 1)
	permissions[0] = "*:*:*"

	var buttons = make([]string, 1)
	buttons[0] = "*:*:*"

	var mp = make(map[string]interface{})
	mp["roles"] = roles
	if rolekey == "admin" || rolekey == "系统管理员" {
		mp["permissions"] = permissions
		mp["buttons"] = buttons
	} else {
		list, _ := models.CRoleMenu.GetPermissionWithRoleId(gcontext.Context(c))
		mp["permissions"] = list
		mp["buttons"] = list
	}

	user, err := models.CUser.GetInfo(gcontext.Context(c))
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}

	mp["introduction"] = " am a super administrator"

	mp["avatar"] = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	if user.Avatar != "" {
		mp["avatar"] = user.Avatar
	}
	mp["userName"] = user.NickName
	mp["userId"] = user.UserId
	mp["deptId"] = user.DeptId
	mp["name"] = user.NickName

	servers.OK(c, servers.WithData(mp))
}
