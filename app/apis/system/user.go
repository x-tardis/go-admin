package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/thinkgos/sharp/core/paginator"

	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/izap"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"
)

type User struct{}

// @Summary 列表用户信息数据
// @Description 获取JSON
// @Tags 用户
// @Param username query string false "username"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/users [get]
// @Security Bearer
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
	servers.JSON(c, http.StatusOK, servers.WithData(&paginator.Pages{
		Info: info,
		List: items,
	}))
}

// @Summary 获取用户
// @Description 获取JSON
// @Tags 用户
// @Param userId path int true "用户编码"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/users/{id} [get]
// @Security Bearer
func (User) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	result, err := models.CUser.Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithPrompt(prompt.NotFound),
			servers.WithError(err))
		return
	}
	roles, err := models.CRole.Query(gcontext.Context(c))
	posts, err := models.CPost.Query(gcontext.Context(c), models.PostQueryParam{})

	postIds := make([]int, 0)
	postIds = append(postIds, result.PostId)

	roleIds := make([]int, 0)
	roleIds = append(roleIds, result.RoleId)
	servers.JSONs(c, http.StatusOK, gin.H{
		"code":    200,
		"data":    result,
		"postIds": postIds,
		"roleIds": roleIds,
		"roles":   roles,
		"posts":   posts,
	})
}

// @Summary 获取个人中心用户
// @Description 获取JSON
// @Tags 个人中心
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/user/profile [get]
// @Security Bearer
func (User) GetProfile(c *gin.Context) {
	result, err := models.CUser.GetUserInfo(gcontext.Context(c))
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
	dept, err := models.CDept.Get(gcontext.Context(c), result.DeptId)

	postIds := []int{result.PostId}
	roleIds := []int{result.RoleId}

	servers.JSONs(c, http.StatusOK, gin.H{
		"code":    200,
		"data":    result,
		"postIds": postIds,
		"roleIds": roleIds,
		"roles":   roles,
		"posts":   posts,
		"dept":    dept,
	})
}

// @Summary 获取用户角色和职位
// @Description 获取JSON
// @Tags 用户
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/users [get]
// @Security Bearer
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
	servers.JSON(c, http.StatusOK, servers.WithData(mp))
}

// @Summary 创建用户
// @Description 获取JSON
// @Tags 用户
// @Accept  application/json
// @Product application/json
// @Param data body models.User true "用户数据"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/users [post]
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
	servers.JSON(c, http.StatusOK, servers.WithPrompt(prompt.CreateSuccess))
}

// @Summary 修改用户数据
// @Description 获取JSON
// @Tags 用户
// @Accept  application/json
// @Product application/json
// @Param data body models.User true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/users [put]
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
	servers.JSON(c, http.StatusOK, servers.WithData(item))
}

// @Summary 删除用户数据
// @Description 删除数据
// @Tags 用户
// @Param userId path int true "userId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/users/{ids} [delete]
func (User) BatchDelete(c *gin.Context) {
	ids := infra.ParseIdsGroup(c.Param("ids"))
	err := models.CUser.BatchDelete(ids)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.DeleteFailed))
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithPrompt(prompt.DeleteSuccess))
}

// @Summary 修改头像
// @Description 获取JSON
// @Tags 用户
// @Accept multipart/form-data
// @Param file formData file true "file"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/user/avatar [post]
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
	servers.JSON(c, http.StatusOK, servers.WithData(avatar))
}

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
	servers.JSON(c, http.StatusOK, servers.WithMsg("密码修改成功"))
}
