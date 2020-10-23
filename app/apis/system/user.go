package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/thinkgos/sharp/core/paginator"

	"github.com/x-tardis/go-admin/models"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/izap"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"
)

// User api user
type User struct{}

// @tags 用户
// @summary 列表用户列表
// @description 列表用户列表
// @security Bearer
// @accept json
// @produce json
// @param username query string false "username"
// @param phone query string false "phone"
// @param status query string false "status"
// @param deptId query string false "deptId"
// @param pageSize query int false "页条数"
// @param pageIndex query int false "页码"
// @success 200 {string} servers.Response	"{"code": 200, "data": [...]}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/users [get]
func (User) QueryPage(c *gin.Context) {
	qp := models.UserQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	qp.Inspect()

	items, info, err := models.CUser.QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithPrompt(prompt.QueryFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(&paginator.Pages{
		Info: info,
		List: items,
	}))
}

// @tags 用户
// @summary 获取用户
// @description 获取用户
// @security Bearer
// @accept json
// @produce json
// @param id path int true "主键"
// @success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/users/{id} [get]
func (User) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	user, err := models.CUser.Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithPrompt(prompt.NotFound),
			servers.WithError(err))
		return
	}
	roles, err := models.CRole.Query(gcontext.Context(c))
	posts, err := models.CPost.Query(gcontext.Context(c), models.PostQueryParam{})

	servers.JSON(c, http.StatusOK, gin.H{
		"code":    200,
		"data":    user,
		"postIds": []int{user.PostId},
		"roleIds": []int{user.RoleId},
		"roles":   roles,
		"posts":   posts,
	})
}

// @tags 用户
// @summary 获取用户角色和职位
// @description 获取用户角色和职位
// @security Bearer
// @accept json
// @produce json
// @success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/users/ [get]
func (User) GetInit(c *gin.Context) {
	roles, err := models.CRole.Query(gcontext.Context(c))
	posts, err := models.CPost.Query(gcontext.Context(c), models.PostQueryParam{})
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithPrompt(prompt.NotFound),
			servers.WithError(err))
		return
	}
	mp := map[string]interface{}{
		"roles": roles,
		"posts": posts,
	}
	servers.OK(c, servers.WithData(mp))
}

// @tags 用户
// @summary 创建用户
// @description 创建用户
// @security Bearer
// @accept json
// @produce json
// @param newItem body models.User true "new item"
// @success 200 {string} string	"{"code": 200, "message": ""}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/users [post]
func (User) Create(c *gin.Context) {
	newItem := models.User{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	_, err := models.CUser.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithPrompt(prompt.CreateSuccess))
}

// @tags 用户
// @summary 修改用户数据
// @description 修改用户数据
// @security Bearer
// @accept json
// @produce json
// @param up body models.User true "update item"
// @success 200 {string} string	"{"code": 200, "msg": ""}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/users [put]
func (User) Update(c *gin.Context) {
	up := models.User{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CUser.Update(gcontext.Context(c), up.UserId, up)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.UpdateFailed))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 用户
// @summary 批量删除用户数据
// @description 批量删除用户数据
// @security Bearer
// @accept json
// @produce json
// @param ids path string true "以','分隔的id列列"
// @success 200 {string} string	"{"code": 200, "msg": ""}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/users/{ids} [delete]
func (User) BatchDelete(c *gin.Context) {
	ids := infra.ParseIdsGroup(c.Param("ids"))
	err := models.CUser.BatchDelete(gcontext.Context(c), ids)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.DeleteFailed))
		return
	}
	servers.OK(c, servers.WithPrompt(prompt.DeleteSuccess))
}

// @tags 个人中心
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
	user, err := models.CUser.GetUserInfo(gcontext.Context(c))
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

// @tags 个人中心
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

// @tags 个人中心
// @summary 修改密码
// @description 修改密码
// @security Bearer
// @accept json
// @produce json
// @param upPwd body models.UpdateUserPwd true "update password"
// @success 200 {string} string	"{"code": 200, "msg": ""}"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/user/avatar [put]
func (User) UpdatePassword(c *gin.Context) {
	upPwd := models.UpdateUserPwd{}
	if err := c.ShouldBindJSON(&upPwd); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	err := models.CUser.UpdatePassword(gcontext.Context(c), upPwd)
	if err != nil {
		servers.Fail(c, http.StatusOK, servers.WithMsg("密码更新失败"))
		return
	}
	servers.OK(c, servers.WithMsg("密码修改成功"))
}
