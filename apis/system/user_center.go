package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinkgos/go-core-package/lib/ternary"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/x-tardis/go-admin/models"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/izap"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"
)

const defaultAvatar = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"

// @tags 个人中心/UserCenter
// @summary 获取用户信息
// @description 获取用户信息
// @security Bearer
// @accept json
// @produce json
// @success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/user/info [get]
func (User) GetInfo(c *gin.Context) {
	user, err := models.CUser.GetInfo(gcontext.Context(c))
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	rolekey := jwtauth.FromRoleKey(gcontext.Context(c))
	mp := map[string]interface{}{
		"roles":       []string{rolekey},
		"permissions": []string{"*:*:*"},
		"buttons":     []string{"*:*:*"},

		"introduction": " am a super administrator",
		"userName":     user.NickName,
		"userId":       user.UserId,
		"deptId":       user.DeptId,
		"name":         user.NickName,
		"avatar":       ternary.IfString(user.Avatar == "", defaultAvatar, user.Avatar),
	}

	if rolekey != models.SuperAdmin {
		list, _ := models.CRoleMenu.GetPermissionWithRoleId(gcontext.Context(c))
		mp["permissions"] = list
		mp["buttons"] = list
	}

	servers.OK(c, servers.WithData(mp))
}

// @tags 个人中心/UserCenter
// @summary 获取个人中心用户信息
// @description 获取个人中心用户信息
// @security Bearer
// @accept json
// @produce json
// @success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/user/profile [get]
func (User) GetProfile(c *gin.Context) {
	user, err := models.CUser.GetViewInfo(gcontext.Context(c))
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithPrompt(prompt.NotFound),
			servers.WithError(err))
		return
	}

	//获取角色列表
	roles, err := models.CRole.Query(gcontext.Context(c))
	//获取职位列表
	posts, err := models.CPost.Query(gcontext.Context(c), models.PostQueryParam{})
	//获取部门列表
	dept, err := models.CDept.Get(gcontext.Context(c), user.DeptId)

	servers.JSON(c, http.StatusOK, gin.H{
		"code":    200,
		"data":    user,
		"postIds": []int{user.PostId},
		"roleIds": []int{user.RoleId},
		"roles":   roles,
		"posts":   posts,
		"dept":    dept,
	})
}

// @tags 个人中心/UserCenter
// @summary 修改头像
// @description 修改头像
// @security Bearer
// @accept multipart/form-data
// @produce json
// @param file formData file true "file"
// @success 200 {string} string	"{"code": 200, "msg": ""}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/user/avatar [post]
func (User) UploadAvatar(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]
	guid := infra.GenerateUUID()
	filPath := "static/uploadfile/" + guid + ".jpg"
	for _, file := range files {
		izap.Sugar.Debug(file.Filename)
		// 上传文件至指定目录
		_ = c.SaveUploadedFile(file, filPath)
	}
	avatar := "/" + filPath
	err := models.CUser.UpdateAvatar(gcontext.Context(c), avatar)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.UpdateFailed))
		return
	}
	servers.OK(c, servers.WithData(avatar))
}

// UpdatePassword 更新用户密码
type UpdatePassword struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// @tags 个人中心/UserCenter
// @summary 修改密码
// @description 修改密码
// @security Bearer
// @accept json
// @produce json
// @param up body UpdatePassword true "update password"
// @success 200 {string} string	"{"code": 200, "msg": ""}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/user/avatar [put]
func (User) UpdatePassword(c *gin.Context) {
	up := UpdatePassword{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	err := models.CUser.UpdatePassword(gcontext.Context(c), up.OldPassword, up.NewPassword)
	if err != nil {
		servers.Fail(c, http.StatusOK, servers.WithMsg("密码更新失败"))
		return
	}
	servers.OK(c, servers.WithMsg("密码修改成功"))
}
