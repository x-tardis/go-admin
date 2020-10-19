package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/thinkgos/sharp/core/paginator"

	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/codes"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/izap"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/servers"
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
	var data models.SysUser

	param := paginator.Param{
		PageIndex: cast.ToInt(c.Query("pageIndex")),
		PageSize:  cast.ToInt(c.Query("pageSize")),
	}
	param.Inspect()

	data.Username = c.Request.FormValue("username")
	data.Status = c.Request.FormValue("status")
	data.Phone = c.Request.FormValue("phone")

	postId := c.Request.FormValue("postId")
	data.PostId = cast.ToInt(postId)

	deptId := c.Request.FormValue("deptId")
	data.DeptId = cast.ToInt(deptId)

	data.DataScope = jwtauth.UserIdStr(c)

	result, pInfo, err := data.GetPage(param)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(&paginator.Pages{
		Info: pInfo,
		List: result,
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
	var SysUser models.SysUser
	SysUser.UserId = cast.ToInt(c.Param("id"))
	result, err := SysUser.Get()
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	roles, err := new(models.CallRole).Query(gcontext.Context(c))
	posts, err := new(models.CallPost).Query(gcontext.Context(c), models.PostQueryParam{})

	postIds := make([]int, 0)
	postIds = append(postIds, result.PostId)

	roleIds := make([]int, 0)
	roleIds = append(roleIds, result.RoleId)
	c.JSON(http.StatusOK, gin.H{
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
	var SysUser models.SysUser

	SysUser.UserId = jwtauth.UserId(c)
	result, err := SysUser.Get()
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}

	//获取角色列表
	roles, err := new(models.CallRole).Query(gcontext.Context(c))
	//获取职位列表
	posts, err := new(models.CallPost).Query(gcontext.Context(c), models.PostQueryParam{})
	//获取部门列表
	dept, err := new(models.CallDept).Get(gcontext.Context(c), result.DeptId)

	postIds := make([]int, 0)
	postIds = append(postIds, result.PostId)

	roleIds := make([]int, 0)
	roleIds = append(roleIds, result.RoleId)

	c.JSON(http.StatusOK, gin.H{
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
	roles, err := new(models.CallRole).Query(gcontext.Context(c))
	posts, err := new(models.CallPost).Query(gcontext.Context(c), models.PostQueryParam{})
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	mp := make(map[string]interface{}, 2)
	mp["roles"] = roles
	mp["posts"] = posts
	servers.OKWithRequestID(c, mp, "")
}

// @Summary 创建用户
// @Description 获取JSON
// @Tags 用户
// @Accept  application/json
// @Product application/json
// @Param data body models.SysUser true "用户数据"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/users [post]
func (User) Create(c *gin.Context) {
	var sysuser models.SysUser

	if err := c.ShouldBindJSON(&sysuser); err != nil {
		servers.Fail(c, 500, codes.DataParseFailed)
		return
	}

	sysuser.CreateBy = jwtauth.UserIdStr(c)
	id, err := sysuser.Insert()
	if err != nil {
		servers.Fail(c, 500, codes.CreatedFail)
		return
	}
	servers.OKWithRequestID(c, id, codes.CreatedSuccess)
}

// @Summary 修改用户数据
// @Description 获取JSON
// @Tags 用户
// @Accept  application/json
// @Product application/json
// @Param data body models.SysUser true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/users [put]
func (User) Update(c *gin.Context) {
	var data models.SysUser
	err := c.Bind(&data)
	if err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}
	data.UpdateBy = jwtauth.UserIdStr(c)
	result, err := data.Update(data.UserId)
	if err != nil {
		servers.Fail(c, -1, codes.UpdatedFail)
		return
	}
	servers.OKWithRequestID(c, result, codes.UpdatedSuccess)
}

// @Summary 删除用户数据
// @Description 删除数据
// @Tags 用户
// @Param userId path int true "userId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/users/{ids} [delete]
func (User) BatchDelete(c *gin.Context) {
	var data models.SysUser

	data.UpdateBy = jwtauth.UserIdStr(c)
	ids := infra.ParseIdsGroup(c.Param("ids"))
	result, err := data.BatchDelete(ids)
	if err != nil {
		servers.Fail(c, 500, codes.DeletedFail)
		return
	}
	servers.OKWithRequestID(c, result, codes.DeletedSuccess)
}

// @Summary 修改头像
// @Description 获取JSON
// @Tags 用户
// @Accept multipart/form-data
// @Param file formData file true "file"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/user/profileAvatar [post]
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
	sysuser := models.SysUser{}
	sysuser.UserId = jwtauth.UserId(c)
	sysuser.Avatar = "/" + filPath
	sysuser.UpdateBy = jwtauth.UserIdStr(c)
	sysuser.Update(sysuser.UserId)
	servers.OKWithRequestID(c, filPath, codes.UpdatedSuccess)
}

func (User) UpdatePassword(c *gin.Context) {
	var pwd models.SysUserPwd
	err := c.Bind(&pwd)
	if err != nil {
		servers.Fail(c, 500, codes.UpdatedFail)
		return
	}
	sysuser := models.SysUser{}
	sysuser.UserId = jwtauth.UserId(c)
	sysuser.SetPwd(pwd)
	servers.OKWithRequestID(c, "", "密码修改成功")
}
